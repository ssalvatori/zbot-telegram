package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var topCommand = TopCommand{}

func TestTopCommandOK(t *testing.T) {
	topCommand.Db = &db.MockZbotDatabase{
		Level:      "7",
		Find_terms: []string{"foo", "bar"},
		Rand_def:   db.DefinitionItem{Term: "foo", Meaning: "bar"},
	}

	result, _ := topCommand.ProcessText("!top", userTest)
	assert.Equal(t, "foo bar", result, "Top Command")

}

func TestTopCommandNotMatch(t *testing.T) {

	result, _ := topCommand.ProcessText("!top6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := topCommand.ProcessText("!top6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestTopCommandError(t *testing.T) {

	topCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}
	_, err := topCommand.ProcessText("!top", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
