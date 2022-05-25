package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xxxsen/tgblock/protos/gen/tgblock"
)

func TestCalcBlockSizeByIndex(t *testing.T) {
	var blockSize int64 = 20 * 1024 * 1024
	var blockCount int64 = 6
	var extSize int64 = 1 * 1024 * 1024
	meta := &tgblock.FileContext{
		Name:     "abc",
		FileSize: int64(blockSize*(blockCount-1) + extSize), //101MB
		FileHash: "abc",
		FileIds:  []string{"1", "2", "3", "4", "5", "6"},
	}
	{
		cnt, err := CalcBlockSizeByIndex(meta, 0)
		assert.NoError(t, err)
		assert.Equal(t, cnt, blockSize)
	}
	{
		cnt, err := CalcBlockSizeByIndex(meta, 1)
		assert.NoError(t, err)
		assert.Equal(t, cnt, blockSize)
	}
	{
		cnt, err := CalcBlockSizeByIndex(meta, 5)
		assert.NoError(t, err)
		assert.Equal(t, cnt, extSize)
	}
}
