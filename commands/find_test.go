package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var findCommand = FindCommand{}

func TestFindCommandOK(t *testing.T) {

	findCommand.Db = &db.MockZbotDatabase{
		Term:    "bar",
		Meaning: "bar",
	}

	var result string

	result, _ = findCommand.ProcessText("!find foo", userTest)
	assert.Equal(t, "bar", result, "Last Command")
	findCommand.Db = &db.MockZbotDatabase{
		Not_found: true,
	}
	result, _ = findCommand.ProcessText("!find lalal", userTest)
	assert.Equal(t, "", result, "Last Command")

}
func TestFindCommandNotMatch(t *testing.T) {

	result, _ := findCommand.ProcessText("!find6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := findCommand.ProcessText("!find6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestFindCommandError(t *testing.T) {

	findCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
		Error:    true,
	}
	_, err := findCommand.ProcessText("!find lala", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
