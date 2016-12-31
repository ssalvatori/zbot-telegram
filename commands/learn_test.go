package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var learnCommand = LearnCommand{}

func TestLearnCommandOK(t *testing.T) {
	learnCommand.Db = &db.MockZbotDatabase{
	}
	assert.Equal(t, "[foo] - [bar]", learnCommand.ProcessText("!learn foo bar", user), "Lean Command")
	assert.Equal(t, "", learnCommand.ProcessText("!learn6", user), "Lean no next command")
	learnCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", learnCommand.ProcessText("??", user), "Lean next command")
}
