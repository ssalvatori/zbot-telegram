package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var ignoreCommand = IgnoreCommand{}

func TestIgnoreCommandHelp(t *testing.T) {
	want := "*!ignore* Options available: \n list (show all user ignored) \n add <username> (ignore a user for 10 minutes)"
	res, _ := ignoreCommand.ProcessText("!ignore help", userTest)
	assert.Equal(t, want, res, "Ignore help")
}

func TestIgnoreCommandList(t *testing.T) {
	ignoreCommand.Db = &db.MockZbotDatabase{
		User_ignored: []db.UserIgnore{
			{Username: "rigo", Since: "12", Until: "22"},
			{Username: "jav", Since: "32", Until: "42"},
		},
	}
	expected := "[ @rigo ] since [01-01-1970 00:00:12 UTC] until [01-01-1970 00:00:22 UTC]/n[ @jav ] since [01-01-1970 00:00:32 UTC] until [01-01-1970 00:00:42 UTC]"
	res, _ := ignoreCommand.ProcessText("!ignore list", userTest)
	assert.Equal(t, expected, res, "Last Command")
}

func TestIgnoreCommandAdd(t *testing.T) {
	ignoreCommand.Db = &db.MockZbotDatabase{
		Level: "1000",
		User_ignored: []db.UserIgnore{
			{Username: "rigo", Since: "12", Until: "12"},
			{Username: "jav", Since: "32", Until: "32"},
		},
	}

	userTest.Level = 1000
	expected := "User [rigo] ignored for 10 minutes"
	res, _ := ignoreCommand.ProcessText("!ignore add rigo", userTest)
	assert.Equal(t, expected, res, "Ignore add Command")

	res, _ = ignoreCommand.ProcessText("!ignore add ssalvatori", userTest)
	assert.Equal(t, "You can't ignore yourself", res, "Ignore add myself")

	ignoreCommand.Db = &db.MockZbotDatabase{
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

	ignoreCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}

	_, err := ignoreCommand.ProcessText("!ignore6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Db error")
}
