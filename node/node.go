package node

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hienduyph/goss/httpx"
	"github.com/hienduyph/goss/logger"
	"golang.org/x/sync/errgroup"

	"github.com/hienduyph/genesis/database"
)

const (
	PollInterval    = 5
	Port            = 8080
	endpointStatus  = "/node/status"
	endpointSync    = "/node/sync"
	endpointAddPeer = "/node/peer"
)

func NewNode(
	db *database.State,
	peerState *PeerState,
	miner *Miner,

	balancesHandler *BalanceHandler,
	txHandler *TxHandler,
	nodeHandler *StateHandler,
	syncHandler *SyncHandler,
	peerHandler *PeerHandler,
) *Node {
	h := chi.NewMux()
	h.Get("/balances/list", httpx.Handle(balancesHandler.List))
	h.Post("/tx/add", httpx.Handle(txHandler.Add))
	h.Get(endpointStatus, httpx.Handle(nodeHandler.Status))
	h.Get(endpointSync, httpx.Handle(syncHandler.FromBlockHandler))
	h.Handle(endpointAddPeer, httpx.Handle(peerHandler.Add))

	addr := fmt.Sprintf(":%v", Port)
	server := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	return &Node{
		server:    server,
		db:        db,
		miner:     miner,
		peerState: peerState,
	}
}

type Node struct {
	server *http.Server

	db        *database.State
	peerState *PeerState
	miner     *Miner
}

func (n *Node) Start(parentCtx context.Context) error {
	eg, ctx := errgroup.WithContext(parentCtx)
	eg.Go(func() error {
		return httpx.Run(ctx, n.server)
	})
	eg.Go(func() error {
		return n.sync(ctx)
	})
	eg.Go(func() error {
		return n.miner.Mine(ctx)
	})
	return eg.Wait()
}

func (n *Node) sync(ctx context.Context) error {
	logger.Info("start the syncing daemon", "peers", n.peerState.knownPeers)
	ticker := time.NewTicker(PollInterval * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Info("[sync] got closed signal")
			return nil

		case <-ticker.C:
			n.peerState.doSync()
		}
	}
}

func (n *Node) Close(ctx context.Context) {
	n.db.Close()
}
