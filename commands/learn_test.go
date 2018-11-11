package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var learnCommand = LearnCommand{}

func TestLearnCommandOK(t *testing.T) {
	learnCommand.Db = &db.MockZbotDatabase{}
	assert.Equal(t, "[foo] - [bar]", learnCommand.ProcessText("!learn foo bar", userTest), "Lean Command")
	assert.Equal(t, "", learnCommand.ProcessText("!learn6", userTest), "Learn no next command")

	learnCommand.Db = &db.MockZbotDatabase{Error: true}
	assert.Equal(t, "", learnCommand.ProcessText("!learn foo bar2", userTest), "Learn Error")

	learnCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", learnCommand.ProcessText("??", userTest), "Learn next command")
}
