package main

import (
	"context"
	"fmt"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/goss/logger"
	"github.com/hienduyph/goss/utils/shutdowns"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Launch the TBB node and its HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, done := shutdowns.NewCtx()
			defer done()

			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			bootstrap := []peer.PeerNode{
				{
					IP:          "18.184.213.146",
					Port:        8080,
					IsBootstrap: true,
					IsActive:    true,
				},
			}

			n, e := newNode(ctx, &database.StateConfig{DataDir: dataDir}, bootstrap)
			logger.FatalIf(e, "create nodes")
			defer n.Close(context.Background())

			fmt.Println("Lauching the TBB node and its API ...")
			if err := n.Start(ctx); err != nil {
				panic(err)
			}
		},
	}
	addDefaultRequiredFlags(cmd)
	return cmd
}
