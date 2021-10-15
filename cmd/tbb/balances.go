package main

import (
	"errors"
	"fmt"

	"github.com/hienduyph/genesis/database"
	"github.com/spf13/cobra"
)

var ErrIncorrectUsage = errors.New("incorrect usage")

func balancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list, ...)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ErrIncorrectUsage
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}
	list := &cobra.Command{
		Use:   "list",
		Short: "Lists all balances",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.NewStateFromDisk()
			if err != nil {
				panic(err)
			}
			defer state.Close()
			fmt.Println("Accounts balances:")
			fmt.Println("___________________")
			fmt.Println("")
			for acc, balance := range state.Balances {
				fmt.Println(fmt.Sprintf("%s: %d", acc, balance))
			}
		},
	}
	cmd.AddCommand(list)
	return cmd
}
