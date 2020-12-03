package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var randCommand = RandCommand{}

func TestRandCommandOK(t *testing.T) {

	randCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
	}

	result, _ := randCommand.ProcessText("!rand", userTest, "testchat", false)
	assert.Equal(t, "[foo] - [bar]", result, "Rand command")

}

func TestRandCommandNotMatch(t *testing.T) {

	result, err := randCommand.ProcessText("!rand6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")
	assert.Equal(t, err, ErrNextCommand, "Command doesn't match")
}

func TestRandCommandError(t *testing.T) {

	randCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}
	_, err := randCommand.ProcessText("!rand", userTest, "testchat", false)
	assert.Error(t, err, "Internal Error")
}
