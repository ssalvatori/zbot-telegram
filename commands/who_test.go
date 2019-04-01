package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var whoCommand = WhoCommand{}

func TestWhoCommand(t *testing.T) {

	whoCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Author:  "ssalvatori",
		Date:    "2017-03-22",
		Level:   "100",
	}
	result, _ := whoCommand.ProcessText("!who foo", userTest)
	assert.Equal(t, "[foo] by [ssalvatori] on [2017-03-22]", result, "Who Command OK")
}

func TestWhoCommandNotMatch(t *testing.T) {

	result, _ := whoCommand.ProcessText("!who6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := whoCommand.ProcessText("!who6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

/*
func TestWhoCommandError(t *testing.T) {

	whoCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}
	_, err := whoCommand.ProcessText("!who", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
*/
