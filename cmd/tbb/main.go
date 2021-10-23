package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
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
	tbbCmd.AddCommand(txCmd())
	tbbCmd.AddCommand(runCmd())

	if err := tbbCmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
