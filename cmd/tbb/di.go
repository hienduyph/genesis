//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node"
)

func newNode(ctx context.Context, stateConfig *database.StateConfig) (*node.Node, error) {
	wire.Build(node.GraphSet)
	return nil, nil
}
