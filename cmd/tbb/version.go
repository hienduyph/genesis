package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Major  = "0"
	Minor  = "1"
	Patch  = "0"
	Verbal = "Tx Add && Balances List"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describe version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s.%s.%s %s", Major, Minor, Patch, Verbal)
	},
}
