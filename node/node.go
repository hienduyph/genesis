package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hienduyph/goss/httpx"
	"golang.org/x/sync/errgroup"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
)

const (
	port           = 8080
	endpointStatus = "/node/status"
	endpointSync   = "/node/sync"
)

func NewNode(
	db *database.State,
	bootstraps []peer.PeerNode,

	balancesHandler *BalanceHandler,
	txHandler *TxHandler,
	nodeHandler *StateHandler,
	syncHandler *SyncHandler,
) *Node {
	h := chi.NewMux()
	h.Get("/balances/list", httpx.Handle(balancesHandler.List))
	h.Post("/tx/add", httpx.Handle(txHandler.Add))
	h.Get(endpointStatus, httpx.Handle(nodeHandler.Status))
	h.Get(endpointSync, httpx.Handle(syncHandler.FromBlockHandler))

	addr := fmt.Sprintf(":%v", port)
	server := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	peers := make(map[string]peer.PeerNode)
	for _, p := range bootstraps {
		peers[p.TcpAddress()] = p
	}
	return &Node{
		server:     server,
		db:         db,
		port:       port,
		knownPeers: peers,
	}
}

type Node struct {
	server *http.Server
	port   uint64

	db *database.State
}

func (n *Node) Start(parentCtx context.Context) error {
	eg, ctx := errgroup.WithContext(parentCtx)
	eg.Go(func() error {
		return httpx.Run(ctx, n.server)
	})
	eg.Go(func() error {
		return n.sync(ctx)
	})
	return eg.Wait()
}

func (n *Node) Close(ctx context.Context) {
	n.db.Close()
}
