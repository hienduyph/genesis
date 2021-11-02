package database

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidBlockHash(t *testing.T) {
	art := assert.New(t)
	hexHash := "000000fa04f816039a4db586086168edfa"
	hash := Hash{}
	hex.Decode(hash[:], []byte(hexHash))
	art.True(IsBlockHashValid(hash))
}
