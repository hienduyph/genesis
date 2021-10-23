package handlers

import (
	"fmt"
	"net/http"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/jsonx"
)

func NewTx(db *database.State) *Tx {
	return &Tx{
		db: db,
	}
}

type Tx struct {
	db *database.State
}

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type TxAddResp struct {
	Hash database.Hash `json:"block_hash"`
}

func (tx *Tx) Add(r *http.Request) (interface{}, error) {
	in := new(TxAddReq)
	if e := jsonx.NewDecoder(r.Body).Decode(in); e != nil {
		return nil, fmt.Errorf("decode body failed: `%s` %w", e.Error(), errorx.ErrBadInput)
	}

	x := database.NewTx(
		database.Account(in.From),
		database.NewAccount(in.To),
		in.Value,
		in.Data,
	)
	if e := tx.db.AddTx(x); e != nil {
		return nil, fmt.Errorf("add tx failed: %w", e)
	}
	hash, err := tx.db.Persist()
	if err != nil {
		return nil, fmt.Errorf("persisted failed: %w", err)
	}
	return &TxAddResp{Hash: hash}, nil
}
