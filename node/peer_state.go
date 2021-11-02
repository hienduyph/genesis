package node

import (
	"fmt"
	"io"
	"net/http"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/goss/jsonx"
	"github.com/hienduyph/goss/logger"
)

func NewPeerState(
	bootstraps []peer.PeerNode,
	advertisingInfo peer.PeerNode,
	db *database.State,
) *PeerState {
	peers := make(map[string]peer.PeerNode)
	for _, p := range bootstraps {
		peers[p.TcpAddress()] = p
	}
	return &PeerState{
		knownPeers: peers,
		ip:         advertisingInfo.IP,
		port:       advertisingInfo.Port,
		db:         db,
	}
}

type PeerState struct {
	knownPeers map[string]peer.PeerNode
	ip         string
	port       uint64
	db         *database.State
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

func (n *PeerState) doSync() {
	logger.Debug("Polling for new peers and status")
	for _, p := range n.knownPeers {
		if n.ip == p.IP && n.port == p.Port {
			continue
		}

		status, e := queryPeerStatus(p)
		if e != nil {
			logger.Error(e, "fetch failed", "peer", p)
			n.RemovePeer(p)
			continue
		}

		if err := n.joinKnownPeers(p); err != nil {
			logger.Error(err, "join knownPeers failed", "peer", p, "status", status)
			continue
		}
		if err := n.syncBlocks(p, status); err != nil {
			logger.Error(err, "syncBlocks failed", "peer", p, "status", status)
			continue
		}

		if err := n.syncKnownPeers(p, status); err != nil {
			logger.Error(err, "sync knownPeers failed", "peer", p, "status", status)
			continue
		}
	}
}

func (ps *PeerState) joinKnownPeers(pe peer.PeerNode) error {
	if pe.Connected {
		logger.Debug("peer connected, early return", "pe", pe)
		return nil
	}
	url := fmt.Sprintf(
		"http://%s%s",
		pe.TcpAddress(),
		AddPeerReq{IP: pe.IP, Port: pe.Port}.AsReqURI(endpointAddPeer),
	)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get req failed: %w", err)
	}

	defer res.Body.Close()
	knownPeer := ps.knownPeers[pe.TcpAddress()]
	knownPeer.Connected = true
	ps.AddPeer(knownPeer)
	return nil
}

func (ps *PeerState) syncBlocks(pe peer.PeerNode, ss *StatusResp) error {
	log := logger.Factory("syncBlocks").WithValues("status", ss)
	if ss.Hash.IsEmpty() {
		log.Info("Peer has no block, skip!")
		return nil
	}
	localBlockNumber := ps.db.LatestBlock().Header.Number
	if ss.Number < localBlockNumber {
		log.Info("Peer has less blocks than us, skip!")
		return nil
	}

	if ss.Number == 0 && !ps.db.LatestBlockHash().IsEmpty() {
		log.Info("This is genesis blocks and we've already had in our chains")
		return nil
	}

	newBlocksCount := ss.Number - localBlockNumber
	if localBlockNumber == 0 && ss.Number == 0 {
		newBlocksCount = 1
	}
	log.Info("founds new blocks from peer", "num", newBlocksCount, "peer", pe.TcpAddress())

	blocks, err := fetchBlocksFromPeer(pe, ps.db.LatestBlockHash())
	if err != nil {
		return fmt.Errorf("fetch blocks from peer failed: %w", err)
	}
	if err := ps.db.AddBlocks(blocks); err != nil {
		return fmt.Errorf("apply block to local failed: %w", err)
	}
	return nil
}

func (ps *PeerState) syncKnownPeers(pe peer.PeerNode, ss *StatusResp) error {
	for _, maybeNewPeer := range ss.KnownPeers {
		if !ps.IsKnownPeer(maybeNewPeer) {
			logger.Debug("Found new peer", "peer", maybeNewPeer.TcpAddress())
			ps.AddPeer(maybeNewPeer)
		}
	}
	return nil
}

func queryPeerStatus(peer peer.PeerNode) (*StatusResp, error) {
	log := logger.Factory("queryPeerStatus")
	url := fmt.Sprintf("http://%s%s", peer.TcpAddress(), endpointStatus)
	log.Info("query peer status", "url", url)
	res, e := http.Get(url)
	if e != nil {
		return nil, fmt.Errorf("fetch failed: %w", e)
	}
	r := new(StatusResp)
	e = readJSON(res, r)
	return r, e
}

func fetchBlocksFromPeer(pe peer.PeerNode, fromBlock database.Hash) ([]database.Block, error) {
	log := logger.Factory("fetchBlocksFromPeer")
	url := fmt.Sprintf("http://%s%s", pe.TcpAddress(), FromBlockReq{fromBlock.Hex()}.AsReqURI(endpointSync))
	log.Info("Importing blocks", "peer", pe.TcpAddress(), "url", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch from blocks failed: %w", err)
	}
	out := new(SyncRes)
	if err := readJSON(res, out); err != nil {
		return nil, fmt.Errorf("read json body failed: %w", err)
	}
	return out.Blocks, nil
}

func readJSON(res *http.Response, dst interface{}) error {
	defer res.Body.Close()
	buf, e := io.ReadAll(res.Body)
	if e != nil {
		return fmt.Errorf("read body failed: %w", e)
	}
	if e := jsonx.Unmarshal(buf, dst); e != nil {
		return fmt.Errorf("read and decode body failed: %w", e)
	}
	return nil
}
