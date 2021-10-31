package node

import (
	"fmt"
	"net/http"

	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/genesis/utils/coders"
	"github.com/hienduyph/goss/errorx"
)

type PeerState struct {
	knownPeers map[string]peer.PeerNode
	ip         string
	port       uint64
}

func (p *PeerState) AddPeer(pe peer.PeerNode) {
	p.knownPeers[pe.TcpAddress()] = pe
}
func (p *PeerState) RemovePeer(pe peer.PeerNode) {
	delete(p.knownPeers, pe.TcpAddress())
}

func (p *PeerState) IsKnownPeer(pe peer.PeerNode) bool {
	if pe.IP == p.ip && pe.Port == p.port {
		return true
	}
	_, isKnowMembers := p.knownPeers[pe.TcpAddress()]
	return isKnowMembers
}

type PeerHandler struct {
	peerState *PeerState
}

var successResp = map[string]interface{}{"success": true}

type AddPeerReq struct {
	IP   string `json:"ip"`
	Port uint64 `json:"port"`
}

func (s *PeerHandler) Add(r *http.Request) (interface{}, error) {
	req := new(AddPeerReq)
	if err := coders.Query.Decode(req, r.URL.Query()); err != nil {
		return nil, fmt.Errorf("input input: `%s`; %w", err.Error(), errorx.ErrBadInput)
	}
	p := peer.PeerNode{
		IP:       req.IP,
		Port:     req.Port,
		IsActive: true,
	}
	s.peerState.AddPeer(p)
	return successResp, nil
}
