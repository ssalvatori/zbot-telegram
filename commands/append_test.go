package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var appendCommand = AppendCommand{}

func TestAppendCommandOK(t *testing.T) {

	appendCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
	}

	result, _ := appendCommand.ProcessText("!append foo bar", userTest)
	assert.Equal(t, "[foo] = [bar]", result, "Append Command")

}

func TestAppendCommandNotMatch(t *testing.T) {

	result, _ := appendCommand.ProcessText("!append6 foor ala", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := appendCommand.ProcessText("!append6 fo lala", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestAppendCommandError(t *testing.T) {

	appendCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
		Error:    true,
	}
	_, err := appendCommand.ProcessText("!append foo lala", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")

	appendCommand.Db = &db.MockZbotDatabase{
		ErrorAppend: true,
	}

	_, err = appendCommand.ProcessText("!append foo bar2", userTest)
	assert.Equal(t, "mock", err.Error(), "Append Error Get")
}
