package node

func createRandomPendingBlock() PendingBlock {
	return NewPendingBlock(
		Hash{},
		0,
		[]Tx{
			NewTx("andrej", "andrej", 3, ""),
			NewTx("andrej", "andrej", 700, "reward"),
		},
	)
}
