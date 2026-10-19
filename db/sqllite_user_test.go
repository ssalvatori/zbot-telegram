package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqliteUserIgnoreList(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	list, err := d.UserIgnoreList()
	assert.NoError(t, err)
	assert.Empty(t, list)

	assert.NoError(t, d.UserIgnoreInsert("user1", 300, "admin", 42))

	list, err = d.UserIgnoreList()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "user1", list[0].Username)
	assert.Equal(t, "admin", list[0].CreatedBy)
	assert.Equal(t, uint(42), list[0].ChannelID)
}

func TestSqliteUserLevel(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()
	level, err := d.UserLevel("user1")
	assert.NoError(t, err)
	assert.Equal(t, "bnil", level)
}

func TestSqliteUserCheckIgnore(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	assert.False(t, d.UserCheckIgnore("user1", 42))

	assert.NoError(t, d.UserIgnoreInsert("user1", 300, "admin", 42))
	assert.True(t, d.UserCheckIgnore("user1", 42))
	assert.True(t, d.UserCheckIgnore("USER1", 42))
	assert.False(t, d.UserCheckIgnore("user1", 43))
}

func TestSqliteUserIgnoreInsert(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()
	assert.NoError(t, d.UserIgnoreInsert("user1", 300, "admin", 42))
	assert.True(t, d.UserCheckIgnore("user1", 42))
}

func TestSqliteUserCleanupIgnorelist(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	assert.NoError(t, d.UserIgnoreInsert("expired", -10, "admin", 42))
	assert.NoError(t, d.UserIgnoreInsert("active", 300, "admin", 42))

	assert.NoError(t, d.UserCleanupIgnorelist())

	list, err := d.UserIgnoreList()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "active", list[0].Username)
}
