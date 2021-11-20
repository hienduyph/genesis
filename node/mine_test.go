package node

import "github.com/hienduyph/genesis/database"

func createRandomPendingBlock() PendingBlock {
	return NewPendingBlock(
		database.Hash{},
		0,
		[]database.SignedTx{
			database.NewTx("andrej", "andrej", 3, ""),
			database.NewTx("andrej", "andrej", 700, "reward"),
		},
	)
}
