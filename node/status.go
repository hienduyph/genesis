package node

import (
	"net/http"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
)

func NewStateHandler(
	db *database.State,
	miner *Miner,

	bootstraps []peer.PeerNode,
) *StateHandler {
	return &StateHandler{
		db:         db,
		miner:      miner,
		bootstraps: bootstraps,
	}
}

type StateHandler struct {
	db         *database.State
	bootstraps []peer.PeerNode
	miner      *Miner
}

type StatusResp struct {
	Hash       database.Hash   `json:"block_hash"`
	Number     uint64          `json:"block_number"`
	KnownPeers []peer.PeerNode `json:"known_peers"`
	PendingTXs []database.Tx   `json:"pending_txs"`
}

func (s *StateHandler) Status(r *http.Request) (interface{}, error) {
	return &StatusResp{
		Hash:       s.db.LatestBlockHash(),
		Number:     s.db.LatestBlock().Header.Number,
		KnownPeers: s.bootstraps,
		PendingTXs: s.miner.getPendingTXsAsArray(),
	}, nil
}
