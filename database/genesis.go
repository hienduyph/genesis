package database

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//go:embed genesis.json
var genesisJson []byte

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

func writeGenesisToDisk(path string) error {
	return ioutil.WriteFile(path, genesisJson, 0644)
}
