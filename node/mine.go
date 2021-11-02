package node

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/logger"
)

func NewPendingBlock()  {
}

type PendingBlock struct {
	parent Hash
	number uint64
	time   uint64
	txs    []Tx
	miner  string
}

func Mine(ctx context.Context, pb PendingBlock) (Block, error) {
	emptyBlock := Block{}
	if len(pb.txs) == 0 {
		return emptyBlock, fmt.Errorf("mining empty blocks is not allowed: %w", errorx.ErrBadInput)
	}
	start := time.Now()
	attempt := 0
	var block Block
	var hash Hash
	var nonce uint64
	var err error
	for !IsBlockHashValid(hash) {
		select {
		case <-ctx.Done():
			return emptyBlock, nil
		default:
		}
		attempt++
		nonce = generateNonce()
		if attempt == 1 || attempt%1000000 == 0 {
			logger.Debug("Mining", "pending", len(pb.txs), "attempt", attempt)
		}
		//block = NewBlock(
		//	pb.parent,
		//	pb.number,
		//	nonce,
		//	pb.time,
		//	pb.miner,
		//	pb.txs,
		//)
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
	)
	return block, nil
}

func generateNonce() uint64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Uint64()
}
