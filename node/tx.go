package node

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/jsonx"
)

func NewTxHandler(db *database.State) *TxHandler {
	return &TxHandler{
		db: db,
	}
}

type TxHandler struct {
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

func (tx *TxHandler) Add(r *http.Request) (interface{}, error) {
	in := new(TxAddReq)
	if e := jsonx.NewDecoder(r.Body).Decode(in); e != nil {
		return nil, fmt.Errorf("decode body failed: `%s` %w", e.Error(), errorx.ErrBadInput)
	}

	t := database.NewTx(
		database.Account(in.From),
		database.NewAccount(in.To),
		in.Value,
		in.Data,
	)
	block := database.NewBlock(
		tx.db.LatestBlockHash(),
		tx.db.NextBlockNumber(),
		uint64(time.Now().Unix()),
		[]database.Tx{t},
	)
	hash, e := tx.db.AddBlock(block)
	if e != nil {
		return nil, fmt.Errorf("add tx failed: %w", e)
	}
	return &TxAddResp{Hash: hash}, nil
}
