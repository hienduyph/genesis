package simple

import (
	"crypto/sha256"
	"time"
)

func NewBlock(
	transactions []string,
	prevHash []byte,
) *Block {
	currentTime := time.Now()
	return &Block{
		timestamp:    currentTime,
		transactions: transactions,
		prevHash:     prevHash,
		Hash:         NewHash(currentTime, transactions, prevHash),
	}
}

type Block struct {
	timestamp    time.Time
	transactions []string
	prevHash     []byte
	Hash         []byte
}

func NewHash(currntTime time.Time, transactions []string, prevHash []byte) []byte {
	input := append([]byte{}, prevHash...)
	input = append(input, currntTime.String()...)
	for _, t := range transactions {
		input = append(input, t...)
	}
	h := sha256.Sum256(input)
	return h[:]
}
