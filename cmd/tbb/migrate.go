package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node"
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

			pendingBlock := node.NewPendingBlock(
				database.Hash{},
				state.NextBlockNumber(),
				database.NewAccount(miner),
				[]database.Tx{
					database.NewTx("q", "q", 3, ""),
					database.NewTx("q", "babayaga", 2000, ""),
					database.NewTx("babayaga", "q", 1, ""),
					database.NewTx("babayaga", "caesar", 1000, ""),
					database.NewTx("babayaga", "q", 50, ""),
				},
			)

			_, err = node.Mine(ctx, pendingBlock)
			logger.FatalIf(err, "mine failed")
		},
	}

	addDefaultRequiredFlags(migrateCmd)
	addMinerFlag(migrateCmd)

	return migrateCmd
}
