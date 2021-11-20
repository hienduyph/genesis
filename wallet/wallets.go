package wallet

import "path/filepath"

const keyStorePath = "keystore"

const (
	Q      = "0x63608270e8ae01Fae8e8a3D1Bb0615B897425C95"
	Baba   = "0x644Ea661b7B93Aee7dA0df1083c28F185C2d7E09"
	Caesar = "0x7349b263275f44c0041c884525e444c2faB3f8EB"
)

func GetKeystoreDirPath(path string) string {
	return filepath.Join(path, keyStorePath)
}
