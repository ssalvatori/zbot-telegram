package command

import (
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ignoreCommand = IgnoreCommand{
	Levels: Levels{Ignore: 1000},
}

func TestIgnoreCommandHelp(t *testing.T) {
	result := "*!ignore* Options available: \n list (show all user ignored) \n add <username> (ignore a user for 10 minutes)"
	assert.Equal(t, result, ignoreCommand.ProcessText("!ignore help", user), "Ignore help")
}

func TestIgnoreCommandList(t *testing.T) {
	ignoreCommand.Db = &db.MockZbotDatabase{
		User_ignored: []db.UserIgnore{
			{Username: "rigo", Since: "12", Until: "22"},
			{Username: "jav", Since: "32", Until: "42"},
		},
	}
	expected := "[ @rigo ] since [01-01-1970 00:00:12 UTC] until [01-01-1970 00:00:22 UTC]/n[ @jav ] since [01-01-1970 00:00:32 UTC] until [01-01-1970 00:00:42 UTC]"
	assert.Equal(t, expected, ignoreCommand.ProcessText("!ignore list", user), "Last Command")
}

func TestIgnoreCommandAdd(t *testing.T) {
	ignoreCommand.Db = &db.MockZbotDatabase{
		Level: "1000",
		User_ignored: []db.UserIgnore{
			{Username: "rigo", Since: "12", Until: "12"},
			{Username: "jav", Since: "32", Until: "32"},
		},
	}
	expected := "User [rigo] ignored for 10 minutes"
	assert.Equal(t, expected, ignoreCommand.ProcessText("!ignore add rigo", user), "Ignore add Command")

	assert.Equal(t, "You can't ignore youself", ignoreCommand.ProcessText("!ignore add ssalvatori", user), "Ignore add myself")

	ignoreCommand.Db = &db.MockZbotDatabase{
		Level: "10",
	}
	assert.Equal(t, "Your level is not enough < 1000", ignoreCommand.ProcessText("!ignore add rigo", user), "Ignore add no enough level")
}

func TestConvertDates(t *testing.T) {

	since := "1488644480"
	until := "1488645080"

	sinceFormated, untilFormated := convertDates(since, until)

	assert.Equal(t, "04-03-2017 16:21:20 UTC", sinceFormated, "format ok")
	assert.Equal(t, "04-03-2017 16:31:20 UTC", untilFormated, "format ok")
}
