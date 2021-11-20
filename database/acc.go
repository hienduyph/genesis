package database

import "github.com/ethereum/go-ethereum/common"

func NewAccount(val string) Account {
	return common.HexToAddress(val)
}

type Account = common.Address
