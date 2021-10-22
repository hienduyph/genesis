package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, e := hex.Decode(h[:], data)
	return e
}

func NewBlock(prevHash Hash, ts uint64, payload []Tx) Block {
	return Block{
		Header: BlockHeader{Parent: prevHash, Time: ts},
		TXs:    payload,
	}
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Time   uint64 `json:"time"`
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func (b Block) Hash() (Hash, error) {
	j, e := json.Marshal(b)
	if e != nil {
		return Hash{}, fmt.Errorf("encode failed: %w", e)
	}
	return sha256.Sum256(j), nil
}
