package database

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const BlockReward = 100

var emptyHash = Hash{}

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.Hex()), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, e := hex.Decode(h[:], data)
	return e
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) IsEmpty() bool {
	return bytes.Equal(emptyHash[:], h[:])
}

func IsBlockHashValid(hash Hash) bool {
	const leading = 3
	for i := 0; i < leading; i++ {
		if fmt.Sprintf("%x", hash[i]) != "0" {
			return false
		}
	}
	return fmt.Sprintf("%x", hash[leading+1]) != "0"
}

func NewBlock(
	prevHash Hash,
	num uint64,
	nonce uint64,
	ts uint64,
	miner Account,
	payload []Tx,
) Block {
	return Block{
		Header: BlockHeader{
			Parent: prevHash,
			Time:   ts,
			Number: num,
			Nonce:  nonce,
			Miner:  miner,
		},
		TXs: payload,
	}
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

type BlockHeader struct {
	Parent Hash    `json:"parent"`
	Time   uint64  `json:"time"`
	Number uint64  `json:"number"`
	Nonce  uint64  `json:"nonce"`
	Miner  Account `json:"miner"`
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
