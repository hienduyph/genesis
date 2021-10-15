package database

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadGenesis(path string) (*genesisSchema, error) {
	buf, e := os.ReadFile(path)
	if e != nil {
		return nil, fmt.Errorf("read genesis failed: %w", e)
	}
	var s genesisSchema
	if e := json.Unmarshal(buf, &s); e != nil {
		return nil, fmt.Errorf("decode genesis failed: %w", e)
	}
	return &s, nil

}

type genesisSchema struct {
	Balances map[Account]uint `json:"balances"`
}
