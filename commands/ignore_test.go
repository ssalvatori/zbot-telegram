package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var ignoreCommand = IgnoreCommand{}

func TestIgnoreCommandHelp(t *testing.T) {
	want := "*!ignore* Options available: \n list (show all user ignored) \n add <username> (ignore a user for 10 minutes)"
	res, _ := ignoreCommand.ProcessText("!ignore help", userTest, "testchat", false)
	assert.Equal(t, want, res, "Ignore help")
}

func TestIgnoreCommandList(t *testing.T) {
	ignoreCommand.Db = &db.ZbotDatabaseMock{
		User_ignored: []db.UserIgnore{
			{Username: "rigo", CreatedAt: 12, ValidUntil: 20, CreatedBy: "admin", Chat: "test_dev"},
			{Username: "jav", CreatedAt: 12, ValidUntil: 20, CreatedBy: "admin", Chat: "test_dev"},
		},
	}
	expected := "[ @rigo ] since [1970-01-01 00:00:12 +0000 UTC] until [1970-01-01 00:00:20 +0000 UTC]/n[ @jav ] since [1970-01-01 00:00:12 +0000 UTC] until [1970-01-01 00:00:20 +0000 UTC]"
	res, _ := ignoreCommand.ProcessText("!ignore list", userTest, "testchat", false)
	assert.Equal(t, expected, res, "Last Command")
}

func TestIgnoreCommandAdd(t *testing.T) {
	ignoreCommand.Db = &db.ZbotDatabaseMock{
		Level: "1000",
		User_ignored: []db.UserIgnore{
			{Username: "rigo", CreatedAt: 12, ValidUntil: 20, CreatedBy: "admin", Chat: "test_dev"},
			{Username: "jav", CreatedAt: 12, ValidUntil: 20, CreatedBy: "admin", Chat: "test_dev"},
		},
	}

	userTest.Level = 1000
	expected := "User [rigo] ignored for 10 minutes"
	res, _ := ignoreCommand.ProcessText("!ignore add rigo", userTest, "testchat", false)
	assert.Equal(t, expected, res, "Ignore add Command")

	res, _ = ignoreCommand.ProcessText("!ignore add ssalvatori", userTest, "testchat", false)
	assert.Equal(t, "You can't ignore yourself", res, "Ignore add myself")

	ignoreCommand.Db = &db.ZbotDatabaseMock{
		Level: "10",
	}

}

func TestConvertDates(t *testing.T) {

	since := "1488644480"
	until := "1488645080"

	sinceFormated, untilFormated := convertDates(since, until)

	assert.Equal(t, "04-03-2017 16:21:20 UTC", sinceFormated, "format ok")
	assert.Equal(t, "04-03-2017 16:31:20 UTC", untilFormated, "format ok")
}

func TestIgnoreError(t *testing.T) {

	ignoreCommand.Db = &db.ZbotDatabaseMock{
		Error: true,
	}

	_, err := ignoreCommand.ProcessText("!ignore6", userTest, "testchat", false)
	assert.Equal(t, "no action in command", err.Error(), "Db error")
}
