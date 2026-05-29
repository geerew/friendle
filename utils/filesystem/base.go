package filesystem

import "github.com/spf13/afero"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// FS wraps an afero filesystem for dependency injection and test doubles
type FS struct {
	afero.Fs
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// New creates a filesystem backed by the given afero instance
func New(backend afero.Fs) *FS {
	return &FS{
		Fs: backend,
	}
}
