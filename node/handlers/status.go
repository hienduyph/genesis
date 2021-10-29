package handlers

import (
	"net/http"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
)

func NewNode(
	db *database.State,

	bootstraps []peer.PeerNode,
) *Node {
	return &Node{
		db:         db,
		bootstraps: bootstraps,
	}
}

type Node struct {
	db         *database.State
	bootstraps []peer.PeerNode
}

type StatusResp struct {
	Hash       database.Hash   `json:"block_hash"`
	Number     uint64          `json:"block_number"`
	KnownPeers []peer.PeerNode `json:"peers_known"`
}

func (s *Node) Status(r *http.Request) (interface{}, error) {
	return &StatusResp{
		Hash:       s.db.LatestBlockHash(),
		Number:     s.db.LatestBlock().Header.Number,
		KnownPeers: s.bootstraps,
	}, nil
}
