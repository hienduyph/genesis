package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	art := assert.New(t)
	priKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	art.NoError(err, "gen prikey")

	pubKey := priKey.PublicKey
	pubKeyBytes := elliptic.Marshal(crypto.S256(), pubKey.X, pubKey.Y)
	pubKeyBytesHash := crypto.Keccak256(pubKeyBytes[1:])

	account := common.BytesToAddress(pubKeyBytesHash[12:])

	msg := []byte("mr q hello the blockchain world!")

	sig, err := Sign(msg, priKey)
	art.NoError(err, "sign failed")

	recoveredPubKey, err := Verify(msg, sig)
	art.NoError(err, "must verify ok")

	recoveredPubKeyBytes := elliptic.Marshal(crypto.S256(), recoveredPubKey.X, recoveredPubKey.Y)
	recoveredPubKeyBytesHash := crypto.Keccak256(recoveredPubKeyBytes[1:])
	recoveredAccount := common.BytesToAddress(recoveredPubKeyBytesHash[12:])
	art.Equal(account.Hex(), recoveredAccount.Hex(), "account recovered must match")
}
