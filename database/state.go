package database

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

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
		txMempool:       make([]Tx, 0, 1000),
		dbFile:          f,
		latestBlockHash: Hash{},
		latestBlock:     Block{},
	}
	for acc, balance := range gen.Balances {
		state.Balances[acc] = balance
	}

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
		if err := state.applyBlock(block.Value); err != nil {
			return nil, fmt.Errorf("apply tx failed: %w", err)
		}
		state.latestBlockHash = block.Key
		state.latestBlock = block.Value
	}
	logger.Info("loaded state", "datadir", dataDir, "balances", state.Balances)
	return state, nil
}

type State struct {
	Balances        map[Account]uint
	txMempool       []Tx
	dbFile          *os.File
	latestBlockHash Hash
	latestBlock     Block
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) LatestBlock() Block {
	return s.latestBlock
}

func (s *State) Close() {
	s.dbFile.Close()
}

func (s *State) AddBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.AddTx(tx); err != nil {
			return fmt.Errorf("add tx failed: %w", err)
		}
	}
	return nil
}

func (s *State) AddTx(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return fmt.Errorf("apply txlog failed: %w", err)
	}
	s.txMempool = append(s.txMempool, tx)
	return nil
}

func (s *State) Persist() (Hash, error) {
	block := NewBlock(
		s.latestBlockHash,
		s.latestBlock.Header.Number+1,
		uint64(time.Now().Unix()),
		s.txMempool,
	)
	logger.Info("prepare", "block", block, "latest", s.latestBlock)
	return s.MigrateBlock(block)
}

// MigrateBlock
// internal funcs, use for migrate only
func (s *State) MigrateBlock(block Block) (Hash, error) {
	logger.Info("prepare", "block", block, "latest", s.latestBlock)
	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, fmt.Errorf("hashed failed: %w", err)
	}
	blockFS := BlockFS{blockHash, block}
	blockFSJSON, err := json.Marshal(blockFS)
	if err != nil {
		return Hash{}, fmt.Errorf("encode blockfs failed: %w", err)
	}

	fmt.Printf("Persist new block to disk: \t %s\n", blockFSJSON)
	if _, err := s.dbFile.Write(append(blockFSJSON, '\n')); err != nil {
		return Hash{}, fmt.Errorf("flush to disk failed: %w", err)
	}
	s.latestBlockHash = blockHash
	s.latestBlock = block

	// reset mem pool
	s.txMempool = s.txMempool[:0]
	return s.latestBlockHash, nil
}

func (s *State) applyBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.apply(tx); err != nil {
			return fmt.Errorf("apply failed: %w", err)
		}
	}
	return nil
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}
	if tx.Value > s.Balances[tx.From] {
		return ErrInsufficientBalance
	}
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}
