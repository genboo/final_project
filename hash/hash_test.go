package hash

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("hash", func(t *testing.T) {
		val := "string for hash"
		hash, err := Make(val)
		require.NotZero(t, hash)
		// require.Equal не принимает uint64, поэтому через форматирование в строку
		require.Equal(t, strconv.FormatUint(hash, 10), "15671524549379031109")
		require.Nil(t, err)
	})
}
