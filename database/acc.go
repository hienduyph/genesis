package database

func NewAccount(val string) Account {
	return Account(val)
}

type Account string
