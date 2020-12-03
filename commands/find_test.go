package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var findCommand = FindCommand{}

func TestFindCommandOK(t *testing.T) {

	findCommand.Db = &db.ZbotDatabaseMock{
		Term:    "bar",
		Meaning: "bar",
	}

	var result string

	result, _ = findCommand.ProcessText("!find foo", userTest, "testchat", false)
	assert.Equal(t, "bar", result, "Last Command")
	findCommand.Db = &db.ZbotDatabaseMock{
		NotFound: true,
	}
	result, _ = findCommand.ProcessText("!find lalal", userTest, "testchat", false)
	assert.Equal(t, "", result, "Last Command")

}
func TestFindCommandNotMatch(t *testing.T) {

	result, _ := findCommand.ProcessText("!find6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := findCommand.ProcessText("!find6", userTest, "testchat", false)
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestFindCommandError(t *testing.T) {

	findCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}
	_, err := findCommand.ProcessText("!find lala", userTest, "testchat", false)
	assert.Error(t, err, "DB Error")
}
