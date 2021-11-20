package database

import (
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hienduyph/goss/jsonx"
)

// each tx cost 50 tokens
const TxFee = uint(50)

func NewTx(from Account, to Account, value uint, nonce uint64, msg string) Tx {
	return Tx{from, to, value, msg, uint64(time.Now().UnixMicro()), nonce}
}

type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
	Time  uint64  `json:"time"`
	Nonce uint64  `json:"nonce"`
}

func (t Tx) Hash() (Hash, error) {
	j, e := t.Encode()
	if e != nil {
		return emptyHash, fmt.Errorf("encode tx failed: %w", e)
	}
	return sha256.Sum256(j), nil
}

func (t Tx) Encode() ([]byte, error) {
	return jsonx.Marshal(t)
}

func NewSignedTx(tx Tx, sign []byte) SignedTx {
	return SignedTx{tx, sign}
}

type SignedTx struct {
	Tx  `json:"tx"`
	Sig []byte `json:"signature"`
}

func (st SignedTx) Hash() (Hash, error) {
	txJSON, err := st.Encode()
	if err != nil {
		return emptyHash, err
	}
	return sha256.Sum256(txJSON), nil
}

func (st SignedTx) IsAuthentic() (bool, error) {
	txHash, err := st.Tx.Hash()
	if err != nil {
		return false, err
	}

	recoveredPubKey, err := crypto.SigToPub(txHash[:], st.Sig)
	if err != nil {
		return false, err
	}

	recoveredPubKeyBytes := elliptic.Marshal(crypto.S256(), recoveredPubKey.X, recoveredPubKey.Y)
	recoveredPubKeyBytesHash := crypto.Keccak256(recoveredPubKeyBytes[1:])
	recoveredAccount := common.BytesToAddress(recoveredPubKeyBytesHash[12:])
	return recoveredAccount.Hex() == st.From.Hex(), nil
}
