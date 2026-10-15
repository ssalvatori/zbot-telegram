package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var ignoreCommand = IgnoreCommand{AdminLevel: 100}

func TestIgnoreCommandHelp(t *testing.T) {
	want := "*!ignore* Options available: \n list (show all user ignored) \n add <username> (ignore a user for 10 minutes)"
	res, _ := ignoreCommand.ProcessText("!ignore help", userTest, 0, "testchat", false)
	assert.Equal(t, want, res, "Ignore help")
}

func TestIgnoreCommandList(t *testing.T) {
	ignoreCommand.Db = &db.ZbotDatabaseMock{
		UserIgnored: []db.UserIgnore{
			{Username: "rigo", CreatedAt: 12, ValidUntil: 20, CreatedBy: "admin"},
			{Username: "jav", CreatedAt: 12, ValidUntil: 20, CreatedBy: "admin"},
		},
	}
	expected := "[ @rigo ] since [1970-01-01 00:00:12 +0000 UTC] until [1970-01-01 00:00:20 +0000 UTC]/n[ @jav ] since [1970-01-01 00:00:12 +0000 UTC] until [1970-01-01 00:00:20 +0000 UTC]"
	res, _ := ignoreCommand.ProcessText("!ignore list", userTest, 0, "testchat", false)
	assert.Equal(t, expected, res, "Last Command")
}

func TestIgnoreCommandAdd(t *testing.T) {
	ignoreCommand.Db = &db.ZbotDatabaseMock{
		Level: "1000",
		UserIgnored: []db.UserIgnore{
			{Username: "rigo", CreatedAt: 12, ValidUntil: 20, CreatedBy: "admin"},
			{Username: "jav", CreatedAt: 12, ValidUntil: 20, CreatedBy: "admin"},
		},
	}
	ignoreCommand.IgnoreDuration = 600

	userTest.Level = 1000
	expected := "User [rigo] ignored for 10 minutes"
	res, _ := ignoreCommand.ProcessText("!ignore add rigo", userTest, 0, "testchat", false)
	assert.Equal(t, expected, res, "Ignore add Command")

	res, _ = ignoreCommand.ProcessText("!ignore add ssalvatori", userTest, 0, "testchat", false)
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

	_, err := ignoreCommand.ProcessText("!ignore6", userTest, 0, "testchat", false)
	assert.Equal(t, "no action in command", err.Error(), "Db error")

	_, err = ignoreCommand.ProcessText("!ignore6", userTest, 0, "testchat", false)
	assert.Error(t, err, "Private message")

}

func TestIgnorePrivateMessage(t *testing.T) {
	_, err := ignoreCommand.ProcessText("!ignore help", userTest, 0, "testchat", true)
	assert.Equal(t, ErrNextCommand, err, "Private message")
}

func TestIgnoreListError(t *testing.T) {
	ignoreCommand.Db = &db.ZbotDatabaseMock{
		Error: true,
	}
	_, err := ignoreCommand.ProcessText("!ignore list", userTest, 0, "testchat", false)
	assert.Error(t, err, "List DB error")
}

func TestIgnoreCommandNotAdmin(t *testing.T) {
	nonAdmin := userTest
	nonAdmin.Level = 0

	res, err := ignoreCommand.ProcessText("!ignore add rigo", nonAdmin, 0, "testchat", false)
	assert.NoError(t, err, "Ignore add non admin")
	assert.Equal(t, "Only channel admins can use this command, your level is not enough < 100", res, "Ignore add non admin")
}
