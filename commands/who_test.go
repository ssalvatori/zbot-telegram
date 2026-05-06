package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var whoCommand = WhoCommand{}

func TestWhoCommand(t *testing.T) {

	whoCommand.Db = &db.ZbotDatabaseMock{
		Term:     "foo",
		Meaning:  "bar",
		Author:   "ssalvatori",
		Level:    "100",
		UpdateAt: 1604225446,
		CreateAt: 1604225446,
	}
	result, _ := whoCommand.ProcessText("!who foo", userTest, "testchat", false)
	assert.Equal(t, "[foo] by [ssalvatori] on [2020-11-01 10:10:46 +0000 UTC] hits [0]", result, "Who Command OK")
}

func TestWhoCommandDisabledChannel(t *testing.T) {
	DisableLearnChannels = []string{"disabled-chat"}
	defer func() { DisableLearnChannels = nil }()
	_, err := whoCommand.ProcessText("!who foo", userTest, "disabled-chat", false)
	assert.Equal(t, ErrLearnDisabledChannel, err)
}

func TestWhoCommandNotMatch(t *testing.T) {

	result, _ := whoCommand.ProcessText("!who6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := whoCommand.ProcessText("!who6", userTest, "testchat", false)
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestWhoCommandSetDb(t *testing.T) {
	whoCommand.SetDb(&db.ZbotDatabaseMock{})
}

func TestWhoCommandNotFound(t *testing.T) {
	whoCommand.Db = &db.ZbotDatabaseMock{
		NotFound: true,
	}
	_, err := whoCommand.ProcessText("!who foo", userTest, "testchat", false)
	assert.EqualError(t, err, "Definition [foo] not found")
}

func TestWhoCommandError(t *testing.T) {

	whoCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}
	_, err := whoCommand.ProcessText("!who foo", userTest, "testchat", false)
	assert.Error(t, err, "DB Error")

	_, err = whoCommand.ProcessText("!who foo", userTest, "testchat", true)
	assert.Error(t, err, "Private Message")
}
