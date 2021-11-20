package node

import (
	"fmt"
	"net/http"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/jsonx"
)

func NewTxHandler(
	db *database.State,
	miner *Miner,
	peers *PeerState,
) *TxHandler {
	return &TxHandler{
		db:    db,
		miner: miner,
		peers: peers,
	}
}

type TxHandler struct {
	db    *database.State
	miner *Miner
	peers *PeerState
}

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type TxAddResp struct {
}

func (tx *TxHandler) Add(r *http.Request) (interface{}, error) {
	in := new(TxAddReq)
	if e := jsonx.NewDecoder(r.Body).Decode(in); e != nil {
		return nil, fmt.Errorf("decode body failed: `%s` %w", e.Error(), errorx.ErrBadInput)
	}

	t := database.NewTx(
		database.NewAccount(in.From),
		database.NewAccount(in.To),
		in.Value,
		in.Data,
	)
	if err := tx.miner.AddPendingTX(t, tx.peers.Current()); err != nil {
		return nil, fmt.Errorf("add pending failed: %w", err)
	}
	return &TxAddResp{}, nil
}
