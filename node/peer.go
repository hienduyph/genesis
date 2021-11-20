package node

import (
	"fmt"
	"io"
	"net/http"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/jsonx"
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

func (s *PeerHandler) Add(r *http.Request) (interface{}, error) {
	req := new(AddPeerReq)
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w. `%s`", errorx.ErrBadInput, err.Error())
	}
	if err := jsonx.Unmarshal(buf, req); err != nil {
		return nil, fmt.Errorf("input input: `%s`; %w", err.Error(), errorx.ErrBadInput)
	}

	logger.Debug("Receive peer req", "req", req, "addr", r.RemoteAddr)
	p := peer.PeerNode{
		IP:       req.IP,
		Port:     req.Port,
		IsActive: true,
		Account:  req.Miner,
	}
	s.peerState.AddPeerOrUpdate(p)
	return successResp, nil
}
