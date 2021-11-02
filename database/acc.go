package database

const TxReward = "reward"

func NewAccount(val string) Account {
	return Account(val)
}

type Account string
