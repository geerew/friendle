package dao

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_CreateGroup(t *testing.T) {
	// Test successfully inserting a group record
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)

		group := &models.Group{Name: "Test Group", CreatedBy: user.ID}
		require.NoError(t, dao.CreateGroup(ctx, group))
		require.NotEmpty(t, group.ID)
	})

	// Test error due to duplicate name
	t.Run("duplicate name", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)

		require.NoError(t, dao.CreateGroup(ctx, &models.Group{Name: "Unique", CreatedBy: user.ID}))

		err := dao.CreateGroup(ctx, &models.Group{Name: "unique", CreatedBy: user.ID})
		require.ErrorContains(t, err, "UNIQUE constraint failed")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_GetGroup(t *testing.T) {
	// Test successfully retrieving a group
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Detail Group", user.ID)

		record, err := dao.GetGroup(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: group.ID}))
		require.NoError(t, err)
		require.Equal(t, group.ID, record.ID)
	})

	// Test no error when retrieving a non-existent group
	t.Run("not found", func(t *testing.T) {
		dao, ctx := setup(t)

		record, err := dao.GetGroup(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: "missing"}))
		require.NoError(t, err)
		require.Nil(t, record)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_ListGroups(t *testing.T) {
	// Test successfully listing groups
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		createTestGroup(t, dao, ctx, "Listed Group", user.ID)

		groups, err := dao.ListGroups(ctx, NewOptions())
		require.NoError(t, err)
		require.NotEmpty(t, groups)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_ListGroups_nameSearch(t *testing.T) {
	// Test successfully searching groups by name
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		createTestGroup(t, dao, ctx, "Alpha Group", user.ID)
		createTestGroup(t, dao, ctx, "Beta Group", user.ID)

		groups, err := dao.ListGroups(ctx, NewOptions().WithGroupNameSearch("alpha"))
		require.NoError(t, err)
		require.Len(t, groups, 1)
		require.Equal(t, "Alpha Group", groups[0].Name)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_UpdateGroup(t *testing.T) {
	// Test successfully updating a group
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Old Name", user.ID)

		group.Name = "New Name"
		require.NoError(t, dao.UpdateGroup(ctx, group))

		record, err := dao.GetGroup(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: group.ID}))
		require.NoError(t, err)
		require.Equal(t, "New Name", record.Name)
	})

	// Test error due to missing ID
	t.Run("missing id", func(t *testing.T) {
		dao, ctx := setup(t)

		err := dao.UpdateGroup(ctx, &models.Group{Name: "x"})
		require.ErrorIs(t, err, utils.ErrId)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_DeleteGroups(t *testing.T) {
	// Test successfully deleting a group
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Delete Group", user.ID)

		require.NoError(t, dao.DeleteGroups(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: group.ID})))

		record, err := dao.GetGroup(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: group.ID}))
		require.NoError(t, err)
		require.Nil(t, record)
	})

	// Test error due to missing where clause
	t.Run("missing where", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.DeleteGroups(ctx, nil), utils.ErrWhere)
	})
}
