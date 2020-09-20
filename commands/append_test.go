package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var appendCommand = AppendCommand{}

func TestAppendCommandOK(t *testing.T) {

	appendCommand.Db = &db.ZbotDatabaseMock{
		Term:    "foo",
		Meaning: "bar",
	}

	result, _ := appendCommand.ProcessText("!append foo bar", userTest, "testchat")
	assert.Equal(t, "[foo] = [bar]", result, "Append Command")

}

func TestAppendCommandNotMatch(t *testing.T) {

	result, _ := appendCommand.ProcessText("!append6 foor ala", userTest, "testchat")
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := appendCommand.ProcessText("!append6 fo lala", userTest, "testchat")
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestAppendCommandError(t *testing.T) {

	appendCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{
			db.Definition{Term: "foo", Meaning: "bar"},
		},
		Error: true,
	}
	_, err := appendCommand.ProcessText("!append foo lala", userTest, "testchat")
	assert.Error(t, err, "DB Error")

	appendCommand.Db = &db.ZbotDatabaseMock{
		Error: true,
	}

	_, err = appendCommand.ProcessText("!append foo bar2", userTest, "testchat")
	assert.Error(t, err, "Append Error Get")
}
