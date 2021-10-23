package handlers

import (
	"net/http"

	"github.com/hienduyph/genesis/database"
)

func NewBalance(
	db *database.State,
) *Balance {
	return &Balance{
		db: db,
	}
}

type Balance struct {
	db *database.State
}

type BalanceListResp struct {
	Hash     database.Hash             `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
}

func (b *Balance) List(r *http.Request) (interface{}, error) {
	resp := BalanceListResp{
		Hash:     b.db.LatestBlockHash(),
		Balances: b.db.Balances,
	}
	return resp, nil
}
