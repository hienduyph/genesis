//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node"
	"github.com/hienduyph/genesis/node/peer"
)

func newNode(
	ctx context.Context,
	stateConfig *database.StateConfig,
	peers []peer.PeerNode,
	advertisingInfo peer.PeerNode,
) (*node.Node, error) {
	wire.Build(node.GraphSet)
	return nil, nil
}
