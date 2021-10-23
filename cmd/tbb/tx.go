package main

import (
	"fmt"

	"github.com/hienduyph/genesis/database"
	"github.com/spf13/cobra"
)

const (
	flagFrom    = "from"
	flagTo      = "to"
	flagValue   = "value"
	flagData    = "data"
	flagDataDir = "datadir"
)

func txCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Interact with transactions (add, ...)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ErrIncorrectUsage
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}
	add := &cobra.Command{
		Use:   "add",
		Short: "Add new tx to database",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			data, _ := cmd.Flags().GetString(flagData)
			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)
			tx := database.NewTx(fromAcc, toAcc, value, data)
			state, err := database.NewState(&database.StateConfig{DataDir: dataDir})
			if err != nil {
				panic(err)
			}
			defer state.Close()
			if err := state.AddTx(tx); err != nil {
				panic(err)
			}
			if _, err := state.Persist(); err != nil {
				panic(err)
			}
			fmt.Println("Tx successfully added to the ledger.")
		},
	}
	add.Flags().String(flagFrom, "", "From what accoutn to send tokens")
	add.MarkFlagRequired(flagFrom)

	add.Flags().String(flagTo, "", "To what account to send tokens")
	add.MarkFlagRequired(flagTo)

	add.Flags().Uint(flagValue, 0, "How many tokens to send")
	add.MarkFlagRequired(flagValue)
	addDefaultRequiredFlags(add)

	add.Flags().String(flagData, "", "Possible values: 'reward'")

	cmd.AddCommand(add)
	return cmd
}
