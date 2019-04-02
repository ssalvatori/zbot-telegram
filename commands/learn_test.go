package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var learnCommand = LearnCommand{}

func TestLearnCommandOK(t *testing.T) {
	var result string
	learnCommand.Db = &db.MockZbotDatabase{}

	result, _ = learnCommand.ProcessText("!learn foo bar", userTest)
	assert.Equal(t, "[foo] - [bar]", result, "Lean Command")

}

func TestLearnCommandNotMatch(t *testing.T) {

	result, _ := learnCommand.ProcessText("!learn6 foor ala", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := learnCommand.ProcessText("!learn6 fo lala", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestLearnCommandError(t *testing.T) {

	learnCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
		Error:    true,
	}
	_, err := learnCommand.ProcessText("!learn foo lala", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
