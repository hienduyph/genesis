package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/genesis/wallet"
	"github.com/hienduyph/goss/logger"
	"github.com/spf13/cobra"
)

var migrateCmd = func() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrates the blockchain database according to new business rules.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer done()

			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			miner, _ := cmd.Flags().GetString(flagMiner)

			state, err := database.NewState(&database.StateConfig{DataDir: dataDir})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			logger.Info("init state", "db", state)

			// some acc for migrations
			q := database.NewAccount(wallet.Q)
			babayaga := database.NewAccount(wallet.Baba)
			caesar := database.NewAccount(wallet.Caesar)

			curr := peer.PeerNode{Account: database.NewAccount(miner)}
			m := node.NewMiner(state, curr)
			m.AddPendingTX(database.NewTx(q, q, 3, ""), curr)
			m.AddPendingTX(database.NewTx(q, babayaga, 2000, ""), curr)
			m.AddPendingTX(database.NewTx(babayaga, q, 1, ""), curr)
			m.AddPendingTX(database.NewTx(babayaga, caesar, 1000, ""), curr)
			m.AddPendingTX(database.NewTx(babayaga, q, 50, ""), curr)

			err = m.MinePendingTXs(ctx)
			logger.FatalIf(err, "mine failed")
		},
	}

	addDefaultRequiredFlags(migrateCmd)
	addMinerFlag(migrateCmd)

	return migrateCmd
}
