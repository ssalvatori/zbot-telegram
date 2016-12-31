package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var ignoreCommand = IgnoreCommand{}

func TestIgnoreCommandHelp(t *testing.T) {
	result := "*!ignore* Options available: \n list (show all user ignored) add <username> (ignore a user for 10 minutes)"
	assert.Equal(t, result, ignoreCommand.ProcessText("!ignore help", user), "Ignore help")
}

func TestIgnoreCommandList(t *testing.T) {
	ignoreCommand.Db = &db.MockZbotDatabase{
		User_ignored: []db.UserIgnore{
			{Username: "rigo", Since: "12", Until: "12",},
			{Username: "jav", Since: "32", Until: "32",},
		},
	}
	expected := "[ @rigo ] since [12] until [12]/n[ @jav ] since [32] until [32]"
	assert.Equal(t, expected, ignoreCommand.ProcessText("!ignore list", user), "Last Command")
}

func TestIgnoreCommandAdd(t *testing.T) {
	lastCommand.Db = &db.MockZbotDatabase{
	}
	assert.Equal(t, "[foo] - [bar]", ignoreCommand.ProcessText("!ignore list", user), "Last Command")
}
