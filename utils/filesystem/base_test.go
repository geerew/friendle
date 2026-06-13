package filesystem

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test successfully using the wrapped afero filesystem
func TestNew(t *testing.T) {
	fs := New(afero.NewMemMapFs())

	require.NoError(t, fs.MkdirAll("/data", 0755))
	require.NoError(t, afero.WriteFile(fs, "/data/token", []byte("secret"), 0600))

	data, err := afero.ReadFile(fs, "/data/token")
	require.NoError(t, err)
	require.Equal(t, "secret", string(data))
}
