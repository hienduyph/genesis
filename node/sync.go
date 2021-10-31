package node

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/genesis/utils/coders"
	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/jsonx"
	"github.com/hienduyph/goss/logger"
)

func NewSyncHandler(db *database.State) *SyncHandler {
	return &SyncHandler{db: db}
}

type SyncHandler struct {
	db *database.State
}

type FromBlockReq struct {
	FromBlock string `json:"fromBlock"`
}

type SyncRes struct {
	Blocks []*database.Block `json:"blocks"`
}

func (s *SyncHandler) FromBlockHandler(r *http.Request) (interface{}, error) {
	d := new(FromBlockReq)
	if e := coders.Query.Decode(d, r.URL.Query()); e != nil {

		return nil, fmt.Errorf("decode params error: %s, %w", e.Error(), errorx.ErrBadInput)
	}
	hash := database.Hash{}
	if e := hash.UnmarshalText([]byte(d.FromBlock)); e != nil {
		return nil, fmt.Errorf("invalid hash: %s, %w", e.Error(), errorx.ErrBadInput)
	}
	blocks, err := s.db.GetBlockAfter(r.Context(), hash)
	if err != nil {
		return nil, fmt.Errorf("read blocks failed: %w", err)
	}
	return &SyncRes{Blocks: blocks}, nil

}

func (s *Node) sync(ctx context.Context) error {
	logger.Info("start the syncing daemon", "peers", s.knownPeers)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Info("[sync] got closed signal")
			return nil

		case <-ticker.C:
			s.fetchNewBlocksAndPeers()
		}
	}
}
func (s *Node) fetchNewBlocksAndPeers() {
	logger.Debug("Polling for new peers and status")
	for _, p := range s.knownPeers {
		status, e := queryPeerStatus(p)
		if e != nil {
			logger.Error(e, "fetch failed", "peer", p)
			continue
		}
		localBlockNumber := s.db.LatestBlock().Header.Number
		if localBlockNumber < status.Number {
			newBlocksCound := status.Number - localBlockNumber
			logger.Info("founds new blocks from peer", "num", newBlocksCound, "peer", p.TcpAddress())
		}
		// add back to peer nodes
		for _, statusPeer := range status.KnownPeers {
			newPeer, isKnowPeer := s.knownPeers[statusPeer.TcpAddress()]
			if !isKnowPeer {
				logger.Debug("Found new peer", "peer", newPeer.TcpAddress())
			}
			s.knownPeers[statusPeer.TcpAddress()] = newPeer
		}
	}
}

func queryPeerStatus(peer peer.PeerNode) (*StatusResp, error) {
	url := fmt.Sprintf("http://%s%s", peer.TcpAddress(), endpointStatus)
	res, e := http.Get(url)
	if e != nil {
		return nil, fmt.Errorf("fetch failed: %w", e)
	}
	defer res.Body.Close()
	buf, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, fmt.Errorf("read body failed: %w", e)
	}
	r := new(StatusResp)
	if e := jsonx.Unmarshal(buf, r); e != nil {
		return nil, fmt.Errorf("read and decode body failed: %w", e)
	}
	return r, nil
}
