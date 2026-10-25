package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var topCommand = TopCommand{}

func TestTopCommandOK(t *testing.T) {
	topCommand.Db = &db.ZbotDatabaseMock{
		Level:     "7",
		FindTerms: []string{"foo", "bar"},
		RandDef:   []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
	}

	result, _ := topCommand.ProcessText("!top", userTest, 0, "testchat", false, false)
	assert.Equal(t, "foo bar", result, "Top Command")

}

func TestTopCommandDisabledChannel(t *testing.T) {
	DisableLearnChannels = []string{"disabled-chat"}
	defer func() { DisableLearnChannels = nil }()
	_, err := topCommand.ProcessText("!top", userTest, 0, "disabled-chat", false, false)
	assert.Equal(t, ErrLearnDisabledChannel, err)
}

func TestTopCommandNotMatch(t *testing.T) {

	result, _ := topCommand.ProcessText("!top6", userTest, 0, "testchat", false, false)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := topCommand.ProcessText("!top6", userTest, 0, "testchat", false, false)
	assert.Equal(t, "no action in command", err.Error(), "Error output doesn't match")
}

func TestTopCommandPrivate(t *testing.T) {
	_, err := topCommand.ProcessText("!top", userTest, 0, "testchat", true, false)
	assert.Equal(t, ErrNextCommand, err, "Private message")
}

func TestTopCommandSetDb(t *testing.T) {
	topCommand.SetDb(&db.ZbotDatabaseMock{})
}

func TestTopCommandWithLimit(t *testing.T) {
	topCommand.Db = &db.ZbotDatabaseMock{
		FindTerms: []string{"foo"},
	}
	result, _ := topCommand.ProcessText("!top 5", userTest, 0, "testchat", false, false)
	assert.Equal(t, "foo", result, "Top with limit")

	result, _ = topCommand.ProcessText("!top 200", userTest, 0, "testchat", false, false)
	assert.Equal(t, "foo", result, "Top with limit capped at 100")
}

func TestTopCommandError(t *testing.T) {

	topCommand.Db = &db.ZbotDatabaseMock{
		Error: true,
	}
	_, err := topCommand.ProcessText("!top", userTest, 0, "testchat", false, false)
	assert.Equal(t, "internal error, check logs", err.Error(), "Db error")

	_, err = topCommand.ProcessText("!top", userTest, 0, "testchat", false, false)
	assert.Error(t, err, "Private mesage")

}
