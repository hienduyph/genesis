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
	"github.com/hienduyph/genesis/node/handlers"
	"github.com/hienduyph/genesis/node/peer"
)

const port = 8080

func NewNode(
	db *database.State,
	bootstraps []peer.PeerNode,

	balancesHandler *handlers.Balance,
	txHandler *handlers.Tx,
	nodeHandler *handlers.Node,
) *Node {
	h := chi.NewMux()
	h.Get("/balances/list", httpx.Handle(balancesHandler.List))
	h.Post("/tx/add", httpx.Handle(txHandler.Add))
	h.Get("/node/status", httpx.Handle(nodeHandler.Status))
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
	server     *http.Server
	port       uint64
	knownPeers map[string]peer.PeerNode

	db *database.State
}

func (s *Node) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return httpx.Run(ctx, s.server)
	})
	eg.Go(func() error {
		return s.sync(ctx)
	})
	return eg.Wait()
}

func (s *Node) Close(ctx context.Context) {
	s.db.Close()
}

func (s *Node) sync(ctx context.Context) error {
	logger.Info("start the syncing daemon")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			s.fetchNewBlocksAndPeers()
		}
	}
}
func (s *Node) fetchNewBlocksAndPeers() {
	for _, p := range s.knownPeers {
		s, e := queryPeerStatus(p)
		if e != nil {
			logger.Error(e, "fetch failed", "peer", p)
			continue
		}

	}
}
