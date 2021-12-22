package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncDec(t *testing.T) {
	id := "abcthijsidasdas"
	key := "hello this is a test"
	timestamp := time.Now().Add(60 * time.Second).Unix()
	sig, err := CreateSig(id, key, timestamp)
	assert.NoError(t, err)
	ok, err := CheckSig(id, key, sig, timestamp)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)
}
