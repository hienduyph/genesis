package node

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/goss/jsonx"
	"github.com/hienduyph/goss/logger"
)

const ContentTypeJSON = "application/json"

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

// AddPeerOrUpdate adds if not exists
func (ps *PeerState) AddPeerOrUpdate(pe peer.PeerNode) {
	if pe.IP == ps.current.IP && pe.Port == ps.current.Port {
		return
	}

	key := pe.TcpAddress()
	p, isKnowMembers := ps.knownPeers[key]
	if isKnowMembers {
		p.Account = pe.Account
	} else {
		logger.Debug("found new peer", "peer", key)
		p = pe
	}
	ps.knownPeers[key] = p
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
		return nil
	}
	body := AddPeerReq{
		IP:    ps.Current().IP,
		Port:  ps.Current().Port,
		Miner: ps.Current().Account,
	}
	buf, err := jsonx.Marshal(body)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(
		"http://%s%s",
		pe.TcpAddress(),
		endpointAddPeer,
	)
	res, err := http.Post(url, ContentTypeJSON, bytes.NewReader(buf))
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
	if newBlocksCount == 0 {
		return nil
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
		ps.AddPeerOrUpdate(maybeNewPeer)
	}
	return nil
}

func (ps *PeerState) syncPendingTXs(p peer.PeerNode, txs []database.SignedTx) error {
	for _, tx := range txs {
		if err := ps.miner.AddPendingTX(tx, p); err != nil {
			return fmt.Errorf("add pending tx failed: %w", err)
		}
	}
	return nil
}

func queryPeerStatus(peer peer.PeerNode) (*StatusResp, error) {
	url := fmt.Sprintf("http://%s%s", peer.TcpAddress(), endpointStatus)
	res, e := http.Get(url)
	if e != nil {
		return nil, fmt.Errorf("fetch failed: %w", e)
	}
	r := new(StatusResp)
	e = readJSON(res, r)
	return r, e
}

func fetchBlocksFromPeer(pe peer.PeerNode, fromBlock database.Hash) ([]database.Block, error) {
	url := fmt.Sprintf("http://%s%s", pe.TcpAddress(), FromBlockReq{fromBlock.Hex()}.AsReqURI(endpointSync))
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
