package node

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/hienduyph/goss/httpx"
	"github.com/hienduyph/goss/logger"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/handlers"
)

const port = 8080

func NewNode(
	db *database.State,
	balancesHandler *handlers.Balance,
	txHandler *handlers.Tx,
) *Node {
	h := chi.NewMux()
	h.Get("/balances/list", httpx.Handle(balancesHandler.List))
	h.Post("/tx/add", httpx.Handle(txHandler.Add))
	addr := fmt.Sprintf(":%v", port)

	server := httpx.NewServer()
	server.Addr = addr
	server.Handler = h

	return &Node{server: server, db: db}
}

type Node struct {
	server *httpx.Server
	db     *database.State
}

func (s *Node) Start(ctx context.Context) error {
	logger.Info("Listening", "addr", s.server.Addr)
	return s.server.Run(ctx)
}

func (s *Node) Close(ctx context.Context) {
	s.db.Close()
}
