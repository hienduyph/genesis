package node

import (
	"github.com/google/wire"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/handlers"
)

var GraphSet = wire.NewSet(
	NewNode,
	handlers.NewBalance,
	handlers.NewTx,
	database.NewState,
)
