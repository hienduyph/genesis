package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/node"
	"github.com/hienduyph/genesis/node/peer"
	"github.com/hienduyph/goss/logger"
	"github.com/spf13/cobra"
)

const (
	bootstrapFlags       = "bootstraps"
	advertisingInfoFlags = "advertising-address"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Launch the TBB node and its HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer done()

			miner, _ := cmd.Flags().GetString(flagMiner)
			raws, _ := cmd.Flags().GetStringSlice(bootstrapFlags)
			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			bootstrap := make([]peer.PeerNode, 0, len(raws))

			for _, item := range raws {
				p, err := hostPortToPeer(item)
				logger.FatalIf(err, "parse bootrap node error", "item", item)
				p.IsBootstrap = true
				p.IsActive = true
				bootstrap = append(bootstrap, p)
			}

			currNodeAddr, _ := cmd.Flags().GetString(advertisingInfoFlags)
			if currNodeAddr == "" {
				containerHostname, err := os.Hostname()
				logger.FatalIf(err, "get host error")
				currNodeAddr = fmt.Sprintf("%s:%d", containerHostname, node.Port)
			}
			nodeInfo, err := hostPortToPeer(currNodeAddr)
			logger.FatalIf(err, "parse bootrap node error", "item", currNodeAddr)
			nodeInfo.Account = database.NewAccount(miner)

			logger.Info("advertising info", "node", nodeInfo)

			n, e := newNode(ctx, &database.StateConfig{DataDir: dataDir}, bootstrap, nodeInfo)
			logger.FatalIf(e, "create nodes")
			defer n.Close(context.Background())

			logger.Info("Lauching the TBB node and its API ...")
			if err := n.Start(ctx); err != nil {
				panic(err)
			}
		},
	}
	addDefaultRequiredFlags(cmd)
	addMinerFlag(cmd)
	cmd.Flags().StringSlice(bootstrapFlags, nil, "list of bootrap nodes")
	cmd.Flags().String(advertisingInfoFlags, "", "host:port for the advertising nodes, default is current ip and port of system")
	return cmd
}

func hostPortToPeer(item string) (peer.PeerNode, error) {
	host, portRaw, err := net.SplitHostPort(item)
	if err != nil {
		return peer.PeerNode{}, err
	}
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		return peer.PeerNode{}, err
	}
	return peer.PeerNode{
		IP:   host,
		Port: uint64(port),
	}, nil
}
