package node

import (
	"net/http"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
)

func NewStateHandler(
	db *database.State,

	bootstraps []peer.PeerNode,
) *StateHandler {
	return &StateHandler{
		db:         db,
		bootstraps: bootstraps,
	}
}

type StateHandler struct {
	db         *database.State
	bootstraps []peer.PeerNode
}

type StatusResp struct {
	Hash       database.Hash   `json:"block_hash"`
	Number     uint64          `json:"block_number"`
	KnownPeers []peer.PeerNode `json:"peers_known"`
}

func (s *StateHandler) Status(r *http.Request) (interface{}, error) {
	return &StatusResp{
		Hash:       s.db.LatestBlockHash(),
		Number:     s.db.LatestBlock().Header.Number,
		KnownPeers: s.bootstraps,
	}, nil
}
