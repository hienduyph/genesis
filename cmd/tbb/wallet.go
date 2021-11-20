package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/console/prompt"
	"github.com/hienduyph/genesis/wallet"
	"github.com/hienduyph/goss/logger"
	"github.com/spf13/cobra"
)

func walletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "Manage blockchain accounts and keys",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ErrIncorrectUsage
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	newAcc := &cobra.Command{
		Use:   "new-account",
		Short: "Creates a new account with a new set of a elliptic-curve Private and Public Keys",
		Run: func(cmd *cobra.Command, args []string) {
			pwd := getPassPhrase("Please enter a password to encrypt the new wallet: ", true)
			dataDir := getDataDirFromCmd(cmd)
			ks := keystore.NewKeyStore(
				wallet.GetKeystoreDirPath(dataDir),
				keystore.StandardScryptN,
				keystore.StandardScryptP,
			)
			acc, err := ks.NewAccount(pwd)
			logger.FatalIf(err, "create acc failed")
			logger.Info("account created", "hex", acc.Address.Hex())
		},
	}
	addDefaultRequiredFlags(newAcc)

	cmd.AddCommand(newAcc)

	return cmd
}

func getPassPhrase(prefix string, confirmation bool) string {
	fmt.Println(prefix)
	pass, err := prompt.Stdin.PromptPassword("Password: ")
	logger.FatalIf(err, "read password failed")

	if confirmation {
		confirm, err := prompt.Stdin.PromptPassword("Repeat Password: ")
		logger.FatalIf(err, "read repeate password failed")
		if confirm != pass {
			utils.Fatalf("passwords do not match: ")
		}
	}
	return pass
}
