package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var getCommand = GetCommand{}

func TestGetCommandOK(t *testing.T) {

	getCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
	}

	result, _ := getCommand.ProcessText("? foo", userTest)
	assert.Equal(t, "[foo] - [bar]", result, "Last Command")

	getCommand.Db = &db.MockZbotDatabase{
		Not_found: true,
	}

	result, _ = getCommand.ProcessText("? foo2", userTest)
	assert.Equal(t, "[foo2] Not found!", result, "Last no next command")

}

func TestGetCommandNotMatch(t *testing.T) {

	result, _ := getCommand.ProcessText("?6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := getCommand.ProcessText("?6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestGetCommandError(t *testing.T) {

	getCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
		Error:    true,
	}
	_, err := getCommand.ProcessText("? foo", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
