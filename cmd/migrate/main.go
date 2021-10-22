package main

import (
	"time"

	"github.com/hienduyph/genesis/database"
)

func main() {
	state, err := database.NewStateFromDisk()
	if err != nil {
		panic(err)
	}
	defer state.Close()
	block0 := database.NewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("q", "q", 3, ""),
			database.NewTx("q", "q", 700, "reward"),
		},
	)
	if err := state.AddBlock(block0); err != nil {
		panic(err)
	}
	block0Hash, _ := state.Persist()
	block1 := database.NewBlock(
		block0Hash,
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("q", "babayaga", 2000, ""),
			database.NewTx("q", "q", 100, "reward"),
			database.NewTx("babayaga", "q", 1, ""),
			database.NewTx("babayaga", "caesar", 1000, ""),
			database.NewTx("babayaga", "q", 50, ""),
			database.NewTx("q", "q", 600, "reward"),
		},
	)
	if err := state.AddBlock(block1); err != nil {
		panic(err)
	}
	state.Persist()
}
