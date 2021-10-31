package node

import (
	"net/http"

	"github.com/hienduyph/genesis/database"
)

func NewBalanceHandler(
	db *database.State,
) *BalanceHandler {
	return &BalanceHandler{
		db: db,
	}
}

type BalanceHandler struct {
	db *database.State
}

type BalanceListResp struct {
	Hash     database.Hash             `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
}

func (b *BalanceHandler) List(r *http.Request) (interface{}, error) {
	resp := BalanceListResp{
		Hash:     b.db.LatestBlockHash(),
		Balances: b.db.Balances,
	}
	return resp, nil
}
