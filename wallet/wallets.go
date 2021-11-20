package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/goss/errorx"
)

const keyStorePath = "keystore"

const (
	Q      = "0x63608270e8ae01Fae8e8a3D1Bb0615B897425C95"
	Baba   = "0x644Ea661b7B93Aee7dA0df1083c28F185C2d7E09"
	Caesar = "0x7349b263275f44c0041c884525e444c2faB3f8EB"
)

func GetKeystoreDirPath(path string) string {
	return filepath.Join(path, keyStorePath)
}

func NewKeystoreAccount(datadir, password string) (common.Address, error) {
	ks := keystore.NewKeyStore(GetKeystoreDirPath(datadir), keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, fmt.Errorf("gen acc failed: %w", err)
	}
	return acc.Address, nil
}

func SignTxWithKeystoreAccount(tx database.Tx, acc common.Address, pwd, keystoreDir string) (database.SignedTx, error) {
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	ksAcc, err := ks.Find(accounts.Account{Address: acc})
	if err != nil {
		return database.SignedTx{}, fmt.Errorf("find acc in keystore failed: %w", err)
	}
	ksAccJSON, err := ioutil.ReadFile(ksAcc.URL.Path)
	if err != nil {
		return database.SignedTx{}, fmt.Errorf("read file: %w", err)
	}
	key, err := keystore.DecryptKey(ksAccJSON, pwd)
	if err != nil {
		return database.SignedTx{}, fmt.Errorf("decrypt key failed: %w", err)
	}

	signedTx, err := SignTx(tx, key.PrivateKey)
	if err != nil {
		return database.SignedTx{}, fmt.Errorf("signtx failed: %w", err)
	}
	return signedTx, nil
}

func Sign(msg []byte, privKey *ecdsa.PrivateKey) (sig []byte, err error) {
	msgHash := sha256.Sum256(msg)
	sig, err = crypto.Sign(msgHash[:], privKey)
	if err != nil {
		return nil, fmt.Errorf("keccak256 sign failed: %w", err)
	}
	if len(sig) != crypto.SignatureLength {
		return nil, fmt.Errorf("wrong sign for signature: %w. Got %d, want %d", errorx.ErrBadInput, len(sig), crypto.SignatureLength)
	}
	return sig, nil
}

func SignTx(tx database.Tx, privKey *ecdsa.PrivateKey) (database.SignedTx, error) {
	rawTx, err := tx.Encode()
	if err != nil {
		return database.SignedTx{}, fmt.Errorf("encode tx failed: %w", err)
	}
	sig, err := Sign(rawTx, privKey)
	if err != nil {
		return database.SignedTx{}, fmt.Errorf("sign failed: %w", err)
	}
	return database.NewSignedTx(tx, sig), nil
}

func Verify(msg, sig []byte) (*ecdsa.PublicKey, error) {
	msgHash := sha256.Sum256(msg)
	recoveredPubKey, err := crypto.SigToPub(msgHash[:], sig)
	if err != nil {
		return nil, fmt.Errorf("unable to recover pubkey: %w", err)
	}
	return recoveredPubKey, nil
}
