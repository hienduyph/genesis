package main

import (
	"context"
	"net"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/goss/logger"
	"github.com/spf13/cobra"
)

var bootstrapFlags = "bootstraps"

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Launch the TBB node and its HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer done()

			raws, _ := cmd.Flags().GetStringSlice(bootstrapFlags)

			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			bootstrap := make([]peer.PeerNode, 0, len(raws))
			for _, item := range raws {
				host, portRaw, err := net.SplitHostPort(item)
				if err != nil {
					logger.FatalIf(err, "parse bootrap node error", "item", item)
				}
				port, err := strconv.Atoi(portRaw)
				if err != nil {
				}
				p := peer.PeerNode{
					IP:          host,
					Port:        uint64(port),
					IsBootstrap: true,
					IsActive:    true,
				}
				bootstrap = append(bootstrap, p)
			}

			n, e := newNode(ctx, &database.StateConfig{DataDir: dataDir}, bootstrap)
			logger.FatalIf(e, "create nodes")
			defer n.Close(context.Background())

			logger.Info("Lauching the TBB node and its API ...")
			if err := n.Start(ctx); err != nil {
				panic(err)
			}
		},
	}
	addDefaultRequiredFlags(cmd)
	cmd.Flags().StringSlice(bootstrapFlags, nil, "list of bootrap nodes")
	return cmd
}
