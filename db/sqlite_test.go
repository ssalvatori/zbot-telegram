package db

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testChannelID is the Channel.ID used for definitions in these tests.
const testChannelID uint = 1

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
	d := &ZbotDatabaseSqlite{File: filepath.Join(t.TempDir(), "zbot.db")}
	require.NoError(t, d.Init())
	sqlDB, err := d.DB.DB()
	require.NoError(t, err)

	d.Close()
	require.Error(t, sqlDB.Ping())
}

func TestSqliteInitInvalidPath(t *testing.T) {
	d := &ZbotDatabaseSqlite{File: "/nonexistent/path/zbot_test.db"}
	err := d.Init()
	assert.Error(t, err)
}

func TestSqliteStatistics(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()
	stats, err := d.Statistics(testChannelID)
	assert.NoError(t, err)
	assert.Equal(t, "0", stats)
}

func TestGORMLoggerUsesApplicationDebugLevel(t *testing.T) {
	tests := []struct {
		name        string
		level       slog.Level
		wantSQLLogs bool
	}{
		{name: "debug", level: slog.LevelDebug, wantSQLLogs: true},
		{name: "info", level: slog.LevelInfo, wantSQLLogs: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			applicationLogger := slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{Level: test.level}))
			gormLogger := newGORMLogger(applicationLogger)

			gormLogger.Trace(context.Background(), time.Now(), func() (string, int64) {
				return "SELECT 1", 1
			}, nil)

			gotSQLLogs := strings.Contains(output.String(), "SELECT 1")
			if gotSQLLogs != test.wantSQLLogs {
				t.Errorf("SQL log visibility = %t, want %t; output = %q", gotSQLLogs, test.wantSQLLogs, output.String())
			}
		})
	}
}

func TestSqliteSetAndGet(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	term, err := d.Set(Definition{Term: "hello", Meaning: "world", ChannelID: testChannelID})
	assert.NoError(t, err)
	assert.Equal(t, "hello", term)

	// Duplicate term should get a numbered suffix.
	term2, err := d.Set(Definition{Term: "hello", Meaning: "world2", ChannelID: testChannelID})
	assert.NoError(t, err)
	assert.Equal(t, "hello1", term2)

	def, err := d.Get("hello", testChannelID)
	assert.NoError(t, err)
	assert.Equal(t, "hello", def.Term)
	assert.Equal(t, "world", def.Meaning)
	assert.Equal(t, testChannelID, def.ChannelID)
}

func TestSqliteSetSameTermDifferentChannels(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	term, err := d.Set(Definition{Term: "hello", Meaning: "one", ChannelID: 1})
	require.NoError(t, err)
	assert.Equal(t, "hello", term)

	term, err = d.Set(Definition{Term: "hello", Meaning: "two", ChannelID: 2})
	require.NoError(t, err)
	assert.Equal(t, "hello", term, "same term in another channel keeps its name")

	def, err := d.Get("hello", 2)
	require.NoError(t, err)
	assert.Equal(t, "two", def.Meaning)

	_, err = d.Get("hello", 3)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestSqliteSetStopsAfterTermCollisions(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	for attempt := 0; attempt < maxSetAttempts; attempt++ {
		term := "taken"
		if attempt > 0 {
			term = fmt.Sprintf("taken%d", attempt)
		}
		require.NoError(t, d._set(term, Definition{Term: term, ChannelID: testChannelID}))
	}

	_, err := d.Set(Definition{Term: "taken", ChannelID: testChannelID})
	require.ErrorIs(t, err, ErrTermCollisionLimit)
}

func TestSqliteGetNotFound(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, err := d.Get("nonexistent", testChannelID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestSqliteIncreaseHits(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, err := d.Set(Definition{Term: "foo", Meaning: "bar", ChannelID: testChannelID})
	assert.NoError(t, err)

	def, err := d.Get("foo", testChannelID)
	assert.NoError(t, err)

	err = d.IncreaseHits(def.ID)
	assert.NoError(t, err)
}

func TestSqliteFind(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "term1", Meaning: "meaning hello", ChannelID: testChannelID})
	_, _ = d.Set(Definition{Term: "term2", Meaning: "other meaning", ChannelID: testChannelID})

	results, err := d.Find("hello", testChannelID, 10)
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	_, err = d.Find("zzznomatch", testChannelID, 10)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestSqliteSearch(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "searchterm", Meaning: "some meaning", ChannelID: testChannelID})

	results, err := d.Search("search", testChannelID, 10)
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	_, err = d.Search("zzznomatch", testChannelID, 10)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestSqliteLast(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "lastterm", Meaning: "some meaning", ChannelID: testChannelID})

	results, err := d.Last(testChannelID, 5)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestSqliteTop(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "topterm", Meaning: "some meaning", ChannelID: testChannelID})

	results, err := d.Top(testChannelID, 5)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestSqliteRand(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "randterm", Meaning: "some meaning", ChannelID: testChannelID})

	results, err := d.Rand(testChannelID, 1)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestSqliteAppend(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "appendme", Meaning: "original", ChannelID: testChannelID})

	err := d.Append(Definition{Term: "appendme", Meaning: "addition", ChannelID: testChannelID, Author: "user1"})
	assert.NoError(t, err)

	def, _ := d.Get("appendme", testChannelID)
	assert.Contains(t, def.Meaning, "original")
	assert.Contains(t, def.Meaning, "addition")

	err = d.Append(Definition{Term: "notexists", ChannelID: testChannelID})
	assert.Error(t, err)
}

func TestSqliteLock(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "lockme", Meaning: "some meaning", ChannelID: testChannelID})

	err := d.Lock(Definition{Term: "lockme", ChannelID: testChannelID, Author: "user1"})
	assert.NoError(t, err)

	// Lock again should fail.
	err = d.Lock(Definition{Term: "lockme", ChannelID: testChannelID})
	assert.Error(t, err)

	// Lock non-existent term.
	err = d.Lock(Definition{Term: "nonexistent", ChannelID: testChannelID})
	assert.Error(t, err)
}

func TestSqliteAppendLocked(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()

	_, _ = d.Set(Definition{Term: "locked", Meaning: "some meaning", ChannelID: testChannelID})
	_ = d.Lock(Definition{Term: "locked", ChannelID: testChannelID, Author: "user1"})

	err := d.Append(Definition{Term: "locked", Meaning: "more", ChannelID: testChannelID})
	assert.ErrorIs(t, err, ErrLocked)
}

func TestSqliteForget(t *testing.T) {
	d, cleanup := setupTestSqliteDB(t)
	defer cleanup()
	err := d.Forget(Definition{Term: "foo", ChannelID: testChannelID})
	assert.NoError(t, err)
}
