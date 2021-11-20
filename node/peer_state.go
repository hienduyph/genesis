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
	miner *Miner,
) *PeerState {
	peers := make(map[string]peer.PeerNode)
	for _, p := range bootstraps {
		peers[p.TcpAddress()] = p
	}
	return &PeerState{
		knownPeers: peers,
		current:    advertisingInfo,
		db:         db,
		miner:      miner,
	}
}

type PeerState struct {
	knownPeers map[string]peer.PeerNode
	current    peer.PeerNode
	db         *database.State
	miner      *Miner
}

func (ps *PeerState) Current() peer.PeerNode {
	return ps.current
}

// AddPeer adds if not exists
func (ps *PeerState) AddPeer(pe peer.PeerNode) {
	if !ps.IsKnownPeer(pe) {
		key := pe.TcpAddress()
		logger.Debug("Found new peer", "peer", key)
		ps.knownPeers[key] = pe
	}
}

func (ps *PeerState) RemovePeer(pe peer.PeerNode) {
	delete(ps.knownPeers, pe.TcpAddress())
}

func (ps *PeerState) KnownPeers() []peer.PeerNode {
	out := make([]peer.PeerNode, 0, len(ps.knownPeers))
	for _, item := range ps.knownPeers {
		out = append(out, item)
	}
	return out
}
func (ps *PeerState) IsKnownPeer(pe peer.PeerNode) bool {
	if pe.IP == ps.current.IP && pe.Port == ps.current.Port {
		return true
	}
	_, isKnowMembers := ps.knownPeers[pe.TcpAddress()]
	return isKnowMembers
}

func (ps *PeerState) doSync() {
	peers := ps.KnownPeers()
	logger.Debug("Polling for new peers and status", "peers", peers)
	for _, p := range peers {
		if ps.current.IP == p.IP && ps.current.Port == p.Port {
			continue
		}

		status, e := queryPeerStatus(p)
		if e != nil {
			logger.Error(e, "fetch failed", "peer", p)
			ps.RemovePeer(p)
			continue
		}

		if err := ps.joinKnownPeers(p); err != nil {
			logger.Error(err, "join knownPeers failed", "peer", p)
			continue
		}
		if err := ps.syncBlocks(p, status); err != nil {
			logger.Error(err, "syncBlocks failed", "peer", p, "txs", len(status.PendingTXs), "nums", status.Number)
			continue
		}

		if err := ps.syncKnownPeers(p, status); err != nil {
			logger.Error(err, "sync knownPeers failed", "peer", p, "status", status.KnownPeers)
			continue
		}

		if err := ps.syncPendingTXs(p, status.PendingTXs); err != nil {
			logger.Error(err, "sync pending tx failed", "total", len(status.PendingTXs), "details", status.PendingTXs)
			continue
		}
	}
}

func (ps *PeerState) joinKnownPeers(pe peer.PeerNode) error {
	if pe.Connected {
		logger.Debug("peer connected, early return", "pe", pe)
		return nil
	}
	uri := AddPeerReq{
		IP:    ps.Current().IP,
		Port:  ps.Current().Port,
		Miner: ps.Current().Account,
	}.AsReqURI(endpointAddPeer)
	url := fmt.Sprintf(
		"http://%s%s",
		pe.TcpAddress(),
		uri,
	)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get req failed: %w", err)
	}

	defer res.Body.Close()
	knownPeer := ps.knownPeers[pe.TcpAddress()]
	knownPeer.Connected = true
	ps.knownPeers[pe.TcpAddress()] = knownPeer
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
	for _, block := range blocks {
		if _, err := ps.db.AddBlock(block); err != nil {
			return fmt.Errorf("addblock to local failed: %w", err)
		}
		ps.miner.newSyncedBlocks <- block
	}
	return nil
}

func (ps *PeerState) syncKnownPeers(pe peer.PeerNode, ss *StatusResp) error {
	for _, maybeNewPeer := range ss.KnownPeers {
		ps.AddPeer(maybeNewPeer)
	}
	return nil
}

func (ps *PeerState) syncPendingTXs(p peer.PeerNode, txs []database.Tx) error {
	for _, tx := range txs {
		if err := ps.miner.AddPendingTX(tx, p); err != nil {
			return fmt.Errorf("add pending tx failed: %w", err)
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
