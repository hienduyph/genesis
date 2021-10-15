package database

func NewTx(from Account, to Account, value uint, msg string) Tx {
	return Tx{from, to, value, msg}
}
