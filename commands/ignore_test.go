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
	result := "*!ignore* Options available: \n list (show all user ignored) add <username> (ignore a user for 10 minutes)"
	assert.Equal(t, result, ignoreCommand.ProcessText("!ignore help", user), "Ignore help")
}

func TestIgnoreCommandList(t *testing.T) {
	ignoreCommand.Db = &db.MockZbotDatabase{
		User_ignored: []db.UserIgnore{
			{Username: "rigo", Since: "12", Until: "12"},
			{Username: "jav", Since: "32", Until: "32"},
		},
	}
	expected := "[ @rigo ] since [12] until [12]/n[ @jav ] since [32] until [32]"
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
	assert.Equal(t, "level not enough (minimum 1000 yours 10)", ignoreCommand.ProcessText("!ignore add rigo", user), "Ignore add no enough level")
}
