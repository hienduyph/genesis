package node

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/genesis/utils/coders"
	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/logger"
)

func NewPeerHandler(
	peerState *PeerState,
) *PeerHandler {
	return &PeerHandler{
		peerState: peerState,
	}
}

// PeerHandler handles p2p system
type PeerHandler struct {
	peerState *PeerState
}

var successResp = map[string]interface{}{"success": true}

type AddPeerReq struct {
	IP    string           `json:"ip"`
	Port  uint64           `json:"port"`
	Miner database.Account `json:"miner"`
}

func (r AddPeerReq) AsReqURI(endpoint string) string {
	u := make(url.Values)
	if e := coders.EncodeQuery(r, u); e != nil {
		logger.Error(e, "encode peer req failed")
	}
	return fmt.Sprintf("%s?%s", endpoint, u.Encode())

}

func (s *PeerHandler) Add(r *http.Request) (interface{}, error) {
	req := new(AddPeerReq)
	if err := coders.DecodeQuery(req, r.URL.Query()); err != nil {
		return nil, fmt.Errorf("input input: `%s`; %w", err.Error(), errorx.ErrBadInput)
	}

	logger.Debug("Receive peer req", "req", req, "addr", r.RemoteAddr)
	p := peer.PeerNode{
		IP:       req.IP,
		Port:     req.Port,
		IsActive: true,
		Account:  req.Miner,
	}
	s.peerState.AddPeer(p)
	return successResp, nil
}
