package node

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/logger"
)

const miningIntervalSeconds = 1

func NewPendingBlock(
	parent database.Hash,
	number uint64,
	miner database.Account,
	txs []database.SignedTx,
) PendingBlock {
	return PendingBlock{
		parent: parent,
		number: number,
		txs:    txs,
		time:   uint64(time.Now().Unix()),
		miner:  miner,
	}
}

type PendingBlock struct {
	parent database.Hash
	number uint64
	time   uint64
	txs    []database.SignedTx
	miner  database.Account
}

func NewMiner(
	db *database.State,
	advertisingInfo peer.PeerNode,
) *Miner {
	return &Miner{
		db:          db,
		currentNode: advertisingInfo,

		pendingTxs:      make(map[string]database.SignedTx, 1000),
		archivedTXs:     make(map[string]database.SignedTx, 1000),
		newSyncedBlocks: make(chan database.Block, 10),
		newPendingTXs:   make(chan database.SignedTx, 10000),
		isMining:        false,
	}

}

type Miner struct {
	db          *database.State
	currentNode peer.PeerNode

	pendingTxs      map[string]database.SignedTx
	archivedTXs     map[string]database.SignedTx
	newSyncedBlocks chan database.Block
	newPendingTXs   chan database.SignedTx
	isMining        bool
}

func (m *Miner) AddPendingTX(tx database.SignedTx, node peer.PeerNode) error {
	txHash, err := tx.Hash()
	if err != nil {
		return fmt.Errorf("hash failed: %w", err)
	}
	hex := txHash.Hex()
	_, isArchived := m.archivedTXs[hex]
	_, isAlreadyPending := m.pendingTxs[hex]
	if !isAlreadyPending && !isArchived {
		logger.Debug("Added TX to mempool", "tx", tx, "peer", node.TcpAddress())
		m.pendingTxs[hex] = tx
		m.newPendingTXs <- tx
	}
	return nil
}

func (m *Miner) Mine(ctx context.Context) error {
	var miningCtx context.Context
	var stopCurrentMining context.CancelFunc
	ticker := time.NewTicker(miningIntervalSeconds * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil

		case block, ok := <-m.newSyncedBlocks:
			if !ok {
				continue
			}
			// stop current mining
			if m.isMining {
				logger.Info("Mining is in progress. Stop infavor of new block")
				blockHas, _ := block.Hash()
				logger.Info("peer mined faster", "block", blockHas.Hex())
				m.removeMinedPendingTXs(block)
				stopCurrentMining()
			}

		case <-ticker.C:
			go func() {
				if len(m.pendingTxs) == 0 || m.isMining {
					return
				}
				logger.Info("Start mining with", "txs", len(m.pendingTxs))
				m.isMining = true
				miningCtx, stopCurrentMining = context.WithCancel(ctx)
				err := m.MinePendingTXs(miningCtx)
				if err != nil {
					logger.Error(err, "ming block failed")
				}
				m.isMining = false
			}()
		}
	}
}

func (m *Miner) removeMinedPendingTXs(block database.Block) {
	archived := make([]string, 0, len(block.TXs))
	for _, tx := range block.TXs {
		txHash, _ := tx.Hash()
		h := txHash.Hex()
		if _, exists := m.pendingTxs[h]; exists {
			m.archivedTXs[h] = tx
			delete(m.pendingTxs, h)
			archived = append(archived, h)
		}
	}
	if len(archived) > 0 {
		logger.Info("Remove Mined pending tx pool", "txs", archived)
	}
}

func (m *Miner) MinePendingTXs(ctx context.Context) error {
	blockToMine := NewPendingBlock(
		m.db.LatestBlockHash(),
		m.db.NextBlockNumber(),
		m.currentNode.Account,
		m.getPendingTXsAsArray(),
	)
	logger.Info("Start to mining", "acc", m.currentNode.Account, "latest", m.db.LatestBlockHash())
	minedBlock, err := Mine(ctx, blockToMine)
	if err != nil {
		return fmt.Errorf("mined block failed: %w", err)
	}
	m.removeMinedPendingTXs(minedBlock)
	_, err = m.db.AddBlock(minedBlock)
	if err != nil {
		return fmt.Errorf("add block failed: %w", err)
	}
	return nil
}

func (m *Miner) getPendingTXsAsArray() []database.SignedTx {
	txs := make([]database.SignedTx, 0, len(m.pendingTxs))
	for _, tx := range m.pendingTxs {
		txs = append(txs, tx)
	}
	return txs
}

func Mine(ctx context.Context, pb PendingBlock) (database.Block, error) {
	emptyBlock := database.Block{}
	if len(pb.txs) == 0 {
		return emptyBlock, fmt.Errorf("mining empty blocks is not allowed: %w", errorx.ErrBadInput)
	}
	start := time.Now()
	attempt := 0
	var block database.Block
	var hash database.Hash
	var nonce uint32
	var err error
	for !database.IsBlockHashValid(hash) {
		select {
		case <-ctx.Done():
			return emptyBlock, ctx.Err()
		default:
		}
		attempt++
		nonce = generateNonce()
		if attempt == 1 || attempt%1000000 == 0 {
			logger.Debug("Mining", "pending", len(pb.txs), "attempt", attempt, "new_nonce", nonce, "miner", pb.miner)
		}
		block = database.NewBlock(
			pb.parent,
			pb.number,
			nonce,
			pb.time,
			pb.miner,
			pb.txs,
		)
		hash, err = block.Hash()
		if err != nil {
			return emptyBlock, fmt.Errorf("hash failed: %w", err)
		}
	}
	logger.Info(
		"Mined Sucessful!",
		"block", block,
		"time", time.Since(start).Seconds(),
		"attempt", attempt,
		"hash", hash.Hex(),
		"nonce", nonce,
		"miner", pb.miner,
	)
	return block, nil
}

func generateNonce() uint32 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Uint32()
}
