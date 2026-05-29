package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test successfully wrapping a Database in a manager
func TestNewManager(t *testing.T) {
	db := &sqliteDB{}
	manager := NewManager(db)

	require.Same(t, db, manager.DataDb)
}
