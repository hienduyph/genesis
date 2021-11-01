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
var emptyHash = Hash{}

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
		if err := applyTXs(block.Value.TXs, state); err != nil {
			return nil, fmt.Errorf("apply tx failed: %w", err)
		}
		state.latestBlockHash = block.Key
		state.latestBlock = block.Value
		state.hasBlock = true
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

func (s *State) AddBlock(b Block) (Hash, error) {
	pendingState := s.copy()
	if err := applyBlock(b, pendingState); err != nil {
		return emptyHash, fmt.Errorf("apply block failed: %w", err)
	}
	blockHash, err := b.Hash()
	if err != nil {
		return emptyHash, fmt.Errorf("gen hash failed: %w", err)
	}
	blockFs := BlockFS{blockHash, b}
	blockFSJSON, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, fmt.Errorf("encode blockfs failed: %w", err)
	}

	fmt.Printf("Persist new block to disk: \t %s\n", blockFSJSON)
	if _, err := s.dbFile.Write(append(blockFSJSON, '\n')); err != nil {
		return Hash{}, fmt.Errorf("flush to disk failed: %w", err)
	}

	s.Balances = pendingState.Balances
	s.latestBlockHash = blockHash
	s.latestBlock = b
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

func (s *State) copy() *State {
	c := &State{
		latestBlock:     s.latestBlock,
		hasBlock:        s.hasBlock,
		latestBlockHash: s.latestBlockHash,
		txMempool:       append([]Tx(nil), s.txMempool...),
		Balances:        make(map[Account]uint, len(s.Balances)),
	}
	for k, v := range s.Balances {
		c.Balances[k] = v
	}
	return c
}

func applyBlock(b Block, s *State) error {
	nextExpectedBlockNumver := s.latestBlock.Header.Number + 1
	if s.hasBlock && b.Header.Number != nextExpectedBlockNumver {
		return fmt.Errorf("next expected block must `%d` not `%d`, %w", nextExpectedBlockNumver, b.Header.Number, errorx.ErrBadInput)
	}

	if s.hasBlock && s.latestBlock.Header.Number > 0 && !reflect.DeepEqual(b.Header.Parent, s.latestBlockHash) {
		return fmt.Errorf("next block parent hash must be `%x` not `%x`", s.latestBlockHash, b.Header.Parent)
	}
	return applyTXs(b.TXs, s)
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
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}
	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("wrong TX. Sender '%s' balance is %d TBB. Tx cost is %d TBB", tx.From, s.Balances[tx.From], tx.Value)
	}
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}
