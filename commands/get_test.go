package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var getCommand = GetCommand{}

func TestGetCommandOK(t *testing.T) {

	getCommand.Db = &db.ZbotDatabaseMock{
		Term:    "foo",
		Meaning: "bar",
	}

	result, _ := getCommand.ProcessText("? foo", userTest, "testchat")
	assert.Equal(t, "[foo] - [bar]", result, "Get Command")

}

func TestGetCommandNoFound(t *testing.T) {
	getCommand.Db = &db.ZbotDatabaseMock{
		NotFound: true,
	}

	result, _ := getCommand.ProcessText("? foo2", userTest, "testchat")
	assert.Equal(t, "[foo2] Not found!", result, "Get no next command")
}

func TestGetCommandNotMatch(t *testing.T) {

	result, _ := getCommand.ProcessText("?6", userTest, "testchat")
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := getCommand.ProcessText("?6", userTest, "testchat")
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestGetCommandError(t *testing.T) {

	getCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}
	_, err := getCommand.ProcessText("? foo", userTest, "testchat")
	assert.Error(t, err, "DB error")
}
