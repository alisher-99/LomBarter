package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadVersion(t *testing.T) {
	t.Parallel()

	file, err := os.Create("version.txt")
	require.NoError(t, err)

	t.Cleanup(func() {
		cErr := file.Close()
		require.NoError(t, cErr)

		rErr := os.Remove("version.txt")
		require.NoError(t, rErr)
	})

	_, err = file.WriteString("v1.0.0")
	require.NoError(t, err)

	version := readVersion("version.txt")
	require.Equal(t, "v1.0.0", version)
}

func TestReadVersion_NotExists(t *testing.T) {
	t.Parallel()

	version := readVersion("version2.txt")
	require.Equal(t, "unknown", version)
}
