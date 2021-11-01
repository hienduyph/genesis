package node

import (
	"github.com/google/wire"

	"github.com/hienduyph/genesis/database"
)

var GraphSet = wire.NewSet(
	NewNode,
	database.NewState,
	NewBalanceHandler,
	NewTxHandler,
	NewStateHandler,
	NewSyncHandler,
	NewPeerState,
	NewPeerHandler,
)
