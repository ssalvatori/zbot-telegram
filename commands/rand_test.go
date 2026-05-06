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

func TestRandCommandWithLimit(t *testing.T) {
	randCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{{Term: "foo", Meaning: "bar"}},
	}

	result, _ := randCommand.ProcessText("!rand 5", userTest, "testchat", false)
	assert.Equal(t, "[foo] - [bar]", result, "Rand with numeric limit")

	result, _ = randCommand.ProcessText("!rand 200", userTest, "testchat", false)
	assert.Equal(t, "[foo] - [bar]", result, "Rand with limit capped at 100")
}

func TestRandCommandDisabledChannel(t *testing.T) {
	DisableLearnChannels = []string{"disabled-chat"}
	defer func() { DisableLearnChannels = nil }()
	_, err := randCommand.ProcessText("!rand", userTest, "disabled-chat", false)
	assert.Equal(t, ErrLearnDisabledChannel, err)
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

	_, err = randCommand.ProcessText("!rand", userTest, "testchat", true)
	assert.Error(t, err, "Private Message")
}
