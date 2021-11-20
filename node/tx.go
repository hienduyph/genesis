package node

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/wallet"
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
	From    string `json:"from"`
	FromPwd string `json:"from_pwd"`
	To      string `json:"to"`
	Value   uint   `json:"value"`
	Data    string `json:"data"`
}

type TxAddResp struct {
	Signed string
}

func (txh *TxHandler) Add(r *http.Request) (interface{}, error) {
	in := new(TxAddReq)
	if e := jsonx.NewDecoder(r.Body).Decode(in); e != nil {
		return nil, fmt.Errorf("decode body failed: `%s` %w", e.Error(), errorx.ErrBadInput)
	}
	from := database.NewAccount(in.From)
	if from.String() == common.HexToAddress("").String() {
		return nil, fmt.Errorf("invalid from: %w", errorx.ErrBadInput)
	}
	if in.FromPwd == "" {
		return nil, fmt.Errorf("missing password: %w", errorx.ErrBadInput)
	}
	tx := database.NewTx(
		from,
		database.NewAccount(in.To),
		in.Value,
		txh.db.GetNextAccountNonce(from),
		in.Data,
	)
	signedtx, err := wallet.SignTxWithKeystoreAccount(tx, from, in.FromPwd, wallet.GetKeystoreDirPath(txh.db.DataDir()))
	if err != nil {
		return nil, err
	}
	if err := txh.miner.AddPendingTX(signedtx, txh.peers.Current()); err != nil {
		return nil, fmt.Errorf("add pending failed: %w", err)
	}
	return &TxAddResp{Signed: base64.StdEncoding.EncodeToString(signedtx.Sig)}, nil
}
