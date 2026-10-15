package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestSqliteDB(t *testing.T) (*ZbotDatabaseSqlite, func()) {
	t.Helper()
	f, err := os.CreateTemp("", "zbot_test_*.db")
	require.NoError(t, err)
	_ = f.Close()
	d := &ZbotDatabaseSqlite{File: f.Name()}
	err = d.Init()
	if err != nil {
		_ = os.Remove(f.Name())
	}
	require.NoError(t, err)
	return d, func() {
		d.Close()
		_ = os.Remove(f.Name())
	}
}

func TestSqliteGetConnectionInfo(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()
	assert.Contains(t, d.GetConnectionInfo(), d.File)
}

func TestSqliteClose(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()
	d.Close() // should not panic
}

func TestSqliteInitInvalidPath(t *testing.T) {
	d := &ZbotDatabaseSqlite{File: "/nonexistent/path/zbot_test.db"}
	err := d.Init()
	assert.Error(t, err)
}

func TestSqliteStatistics(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()
	stats, err := d.Statistics("testchat")
	assert.NoError(t, err)
	assert.Equal(t, "0", stats)
}

func TestSqliteSetAndGet(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	term, err := d.Set(Definition{Term: "hello", Meaning: "world", Chat: "testchat"})
	assert.NoError(t, err)
	assert.Equal(t, "hello", term)

	// Duplicate term should get a numbered suffix.
	term2, err := d.Set(Definition{Term: "hello", Meaning: "world2", Chat: "testchat"})
	assert.NoError(t, err)
	assert.Equal(t, "hello1", term2)

	def, err := d.Get("hello", "testchat")
	assert.NoError(t, err)
	assert.Equal(t, "hello", def.Term)
	assert.Equal(t, "world", def.Meaning)
}

func TestSqliteGetNotFound(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, err := d.Get("nonexistent", "testchat")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestSqliteIncreaseHits(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, err := d.Set(Definition{Term: "foo", Meaning: "bar", Chat: "testchat"})
	assert.NoError(t, err)

	def, err := d.Get("foo", "testchat")
	assert.NoError(t, err)

	err = d.IncreaseHits(def.ID)
	assert.NoError(t, err)
}

func TestSqliteFind(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "term1", Meaning: "meaning hello", Chat: "testchat"})
	_, _ = d.Set(Definition{Term: "term2", Meaning: "other meaning", Chat: "testchat"})

	results, err := d.Find("hello", "testchat", 10)
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	_, err = d.Find("zzznomatch", "testchat", 10)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestSqliteSearch(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "searchterm", Meaning: "some meaning", Chat: "testchat"})

	results, err := d.Search("search", "testchat", 10)
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	_, err = d.Search("zzznomatch", "testchat", 10)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestSqliteLast(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "lastterm", Meaning: "some meaning", Chat: "testchat"})

	results, err := d.Last("testchat", 5)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestSqliteTop(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "topterm", Meaning: "some meaning", Chat: "testchat"})

	results, err := d.Top("testchat", 5)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestSqliteRand(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "randterm", Meaning: "some meaning", Chat: "testchat"})

	results, err := d.Rand("testchat", 1)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestSqliteAppend(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "appendme", Meaning: "original", Chat: "testchat"})

	err := d.Append(Definition{Term: "appendme", Meaning: "addition", Chat: "testchat", Author: "user1"}, "testchat")
	assert.NoError(t, err)

	def, _ := d.Get("appendme", "testchat")
	assert.Contains(t, def.Meaning, "original")
	assert.Contains(t, def.Meaning, "addition")

	err = d.Append(Definition{Term: "notexists", Chat: "testchat"}, "testchat")
	assert.Error(t, err)
}

func TestSqliteLock(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "lockme", Meaning: "some meaning", Chat: "testchat"})

	err := d.Lock(Definition{Term: "lockme", Chat: "testchat", Author: "user1"}, "testchat")
	assert.NoError(t, err)

	// Lock again should fail.
	err = d.Lock(Definition{Term: "lockme", Chat: "testchat"}, "testchat")
	assert.Error(t, err)

	// Lock non-existent term.
	err = d.Lock(Definition{Term: "nonexistent", Chat: "testchat"}, "testchat")
	assert.Error(t, err)
}

func TestSqliteAppendLocked(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "locked", Meaning: "some meaning", Chat: "testchat"})
	_ = d.Lock(Definition{Term: "locked", Chat: "testchat", Author: "user1"}, "testchat")

	err := d.Append(Definition{Term: "locked", Meaning: "more", Chat: "testchat"}, "testchat")
	assert.ErrorIs(t, err, ErrLocked)
}

func TestSqliteForget(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()
	err := d.Forget(Definition{Term: "foo"}, "chat")
	assert.NoError(t, err)
}

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
