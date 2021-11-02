package database

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/hienduyph/goss/jsonx"
)

func NewTx(from Account, to Account, value uint, msg string) Tx {
	return Tx{from, to, value, msg, uint64(time.Now().Unix())}
}

type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
	Time  uint64  `json:"time"`
}

func (t Tx) IsReward() bool {
	return t.Data == TxReward
}

func (t Tx) Hash() (Hash, error) {
	j, e := jsonx.Marshal(t)
	if e != nil {
		return emptyHash, fmt.Errorf("encode tx failed: %w", e)
	}
	return sha256.Sum256(j), nil
}
