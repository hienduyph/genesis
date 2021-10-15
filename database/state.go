package database

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrInsufficientBalance = errors.New("insufficient balance")

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("read cwd failed: %w", err)
	}
	genesisFile := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genesisFile)
	if err != nil {
		return nil, err
	}
	txLogs := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txLogs, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("open tx logs file failed: %w", err)
	}

	// read line by line in logs
	scaner := bufio.NewScanner(f)
	state := &State{Balances: make(map[Account]uint, 1000), txMempool: make([]Tx, 0, 1000), dbFile: f}
	for acc, balance := range gen.Balances {
		state.Balances[acc] = balance
	}

	for scaner.Scan() {
		if err := scaner.Err(); err != nil {
			return nil, fmt.Errorf("scane txlogs failed: %w", err)
		}
		var tx Tx
		if err := json.Unmarshal(scaner.Bytes(), &tx); err != nil {
			return nil, fmt.Errorf("invalid row txlog: %w", err)
		}
		if err := state.apply(tx); err != nil {
			return nil, fmt.Errorf("apply tx failed: %w", err)
		}
	}
	return state, nil
}

type State struct {
	Balances  map[Account]uint
	txMempool []Tx
	dbFile    *os.File
}

func (s *State) Close() {
	s.dbFile.Close()
}

func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return fmt.Errorf("apply txlog failed: %w", err)
	}
	s.txMempool = append(s.txMempool, tx)
	return nil
}

func (s *State) Persist() error {
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for _, tx := range mempool {
		fmt.Printf("Flushing from %s, to %s, total %v\n", tx.From, tx.To, tx.Value)
		txJSON, err := json.Marshal(tx)
		if err != nil {
			return fmt.Errorf("encode transactions failed: %w", err)
		}
		if _, err := s.dbFile.Write(append(txJSON, '\n')); err != nil {
			return fmt.Errorf("flush to disk failed: %w", err)
		}
		s.txMempool = s.txMempool[1:]
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
