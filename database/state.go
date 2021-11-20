package database

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/jsonx"
	"github.com/hienduyph/goss/logger"
)

var ErrInsufficientBalance = errors.New("insufficient balance")

type StateConfig struct {
	DataDir string
}

func NewState(c *StateConfig) (*State, error) {
	dataDir := c.DataDir
	if err := initDataIfNotExists(dataDir); err != nil {
		return nil, fmt.Errorf("init data dir failed: %w", err)
	}
	genesisFile := getGenesisJSONPathFile(dataDir)
	gen, err := loadGenesis(genesisFile)
	if err != nil {
		return nil, err
	}

	txLogs := getBlocksDBFilePath(dataDir)
	f, err := os.OpenFile(txLogs, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("open tx logs file failed: %w", err)
	}

	// read line by line in logs
	scaner := bufio.NewScanner(f)
	state := &State{
		Balances:        make(map[Account]uint, 1000),
		dbFile:          f,
		latestBlockHash: Hash{},
		latestBlock:     Block{},
		conf:            c,
		hasBlock:        false,
	}
	for acc, balance := range gen.Balances {
		state.Balances[acc] = balance
	}
	logger.Debug("initial states", "balances", state.Balances)

	for scaner.Scan() {
		if err := scaner.Err(); err != nil {
			return nil, fmt.Errorf("scane txlogs failed: %w", err)
		}
		b := scaner.Bytes()
		if len(b) == 0 {
			break
		}
		var block BlockFS
		if err := json.Unmarshal(b, &block); err != nil {
			return nil, fmt.Errorf("invalid row txlog: %w", err)
		}
		h, _ := block.Value.Hash()
		logger.Info("applying block", "block", block.Value, "h", h.Hex())
		if err := applyBlock(block.Value, state); err != nil {
			return nil, fmt.Errorf("apply block failed: %w", err)
		}
		state.latestBlockHash = block.Key
		state.latestBlock = block.Value
		state.hasBlock = true
	}
	logger.Info("loaded state", "datadir", dataDir, "balances", state.Balances)
	return state, nil
}

type State struct {
	Balances        map[Account]uint `json:"balances"`
	dbFile          *os.File
	latestBlockHash Hash
	latestBlock     Block
	conf            *StateConfig
	hasBlock        bool
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) LatestBlock() Block {
	return s.latestBlock
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) NextBlockNumber() uint64 {
	if !s.hasBlock {
		return 0
	}
	return s.latestBlock.Header.Number + 1
}

func (s *State) AddBlocks(blocks []Block) error {
	for _, b := range blocks {
		if _, e := s.AddBlock(b); e != nil {
			return e
		}
	}
	return nil
}

func (s *State) AddBlock(block Block) (Hash, error) {
	blockHash, err := block.Hash()
	if err != nil {
		return emptyHash, fmt.Errorf("gen hash failed: %w", err)
	}
	pendingState := s.copy()
	if err := applyBlock(block, pendingState); err != nil {
		return emptyHash, fmt.Errorf("apply block failed: %w", err)
	}
	blockFs := BlockFS{blockHash, block}
	blockFSJSON, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, fmt.Errorf("encode blockfs failed: %w", err)
	}

	logger.Debug("Persist new block to disk", "block", blockHash.Hex())
	if _, err := s.dbFile.Write(append(blockFSJSON, '\n')); err != nil {
		return Hash{}, fmt.Errorf("flush to disk failed: %w", err)
	}

	s.Balances = pendingState.Balances
	s.latestBlockHash = blockHash
	s.latestBlock = block
	s.hasBlock = true
	return blockHash, nil
}

func (s *State) GetBlockAfter(ctx context.Context, hash Hash) ([]Block, error) {
	f, e := os.OpenFile(getBlocksDBFilePath(s.conf.DataDir), os.O_RDONLY, 0600)
	if e != nil {
		return nil, fmt.Errorf("read dbfile faild: %w", e)
	}
	defer f.Close()
	blocks := make([]Block, 0, 1000)
	// compare initial hash
	shouldStartCollect := reflect.DeepEqual(hash, Hash{})
	scaner := bufio.NewScanner(f)
	for scaner.Scan() {
		if err := scaner.Err(); err != nil {
			return nil, err
		}
		blockFs := new(BlockFS)
		if err := jsonx.Unmarshal(scaner.Bytes(), blockFs); err != nil {
			return nil, err
		}
		if shouldStartCollect {
			blocks = append(blocks, blockFs.Value)
			continue
		}
		shouldStartCollect = blockFs.Key == hash
	}
	return blocks, nil
}

func (s *State) apply(tx Tx) error {
	if tx.Value > s.Balances[tx.From] {
		return ErrInsufficientBalance
	}
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}

func (s *State) copy() *State {
	c := &State{
		latestBlock:     s.latestBlock,
		hasBlock:        s.hasBlock,
		latestBlockHash: s.latestBlockHash,
		Balances:        make(map[Account]uint, len(s.Balances)),
	}
	for k, v := range s.Balances {
		c.Balances[k] = v
	}
	return c
}

func applyBlock(b Block, s *State) error {
	// validate the hash
	hash, err := b.Hash()
	if err != nil {
		return fmt.Errorf("hash failed: %w", err)
	}
	if !IsBlockHashValid(hash) {
		return fmt.Errorf("invalid block hash:`%x`; %w,", hash, errorx.ErrBadInput)
	}

	nextExpectedBlockNumver := s.latestBlock.Header.Number + 1
	if s.hasBlock && b.Header.Number != nextExpectedBlockNumver {
		return fmt.Errorf("next expected block must `%d` not `%d`, %w", nextExpectedBlockNumver, b.Header.Number, errorx.ErrBadInput)
	}

	if s.hasBlock && s.latestBlock.Header.Number > 0 && !reflect.DeepEqual(b.Header.Parent, s.latestBlockHash) {
		return fmt.Errorf("next block parent hash must be `%x` not `%x`", s.latestBlockHash, b.Header.Parent)
	}

	if err := applyTXs(b.TXs, s); err != nil {
		return fmt.Errorf("apply txs failed: %w", err)
	}

	// reward for miner
	s.Balances[b.Header.Miner] += BlockReward
	return nil
}

func applyTXs(txs []Tx, s *State) error {
	for _, tx := range txs {
		if err := applyTx(tx, s); err != nil {
			return err
		}
	}
	return nil
}

func applyTx(tx Tx, s *State) error {
	logger.Debug("applyTx", "tx", tx, "balances", s.Balances)
	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("wrong TX. Sender '%s' balance is %d TBB. Tx cost is %d TBB", tx.From, s.Balances[tx.From], tx.Value)
	}
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}
