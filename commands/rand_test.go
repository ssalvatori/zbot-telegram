package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var randCommand = RandCommand{}

func TestRandCommandOK(t *testing.T) {

	randCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
	}

	result, _ := randCommand.ProcessText("!rand", userTest)
	assert.Equal(t, "[foo] - [bar]", result, "Rand command")

}

func TestRandCommandNotMatch(t *testing.T) {

	result, _ := randCommand.ProcessText("!rand6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := randCommand.ProcessText("!rand6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestRandCommandError(t *testing.T) {

	randCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
		Error:    true,
	}
	_, err := randCommand.ProcessText("!rand", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
