package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	flagFrom    = "from"
	flagTo      = "to"
	flagValue   = "value"
	flagData    = "data"
	flagDataDir = "datadir"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	tbbCmd := &cobra.Command{
		Use:   "tbb",
		Short: "The BlockChain Bar CLI",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	tbbCmd.AddCommand(runCmd())
	tbbCmd.AddCommand(migrateCmd())

	if err := tbbCmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
