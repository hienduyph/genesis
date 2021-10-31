package node

import (
	"github.com/google/wire"

	"github.com/hienduyph/genesis/database"
)

var GraphSet = wire.NewSet(
	NewBalanceHandler,
	NewTxHandler,
	NewStateHandler,
	NewNode,
	database.NewState,
)
