package main

import "github.com/spf13/cobra"

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to the node data dir where the DB will/is stored")
	cmd.MarkFlagRequired(flagDataDir)
}

func addMinerFlag(cmd *cobra.Command) {
	cmd.Flags().String(flagMiner, "", "name for this miner")
	cmd.MarkFlagRequired(flagMiner)
}
