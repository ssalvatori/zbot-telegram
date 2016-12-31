package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var levelCommand = LevelCommand{}

func TestLevelCommandOK(t *testing.T) {
	levelCommand.Db = &db.MockZbotDatabase{
		Level: "1000",
	}
	assert.Equal(t, "ssalvatori level 1000", levelCommand.ProcessText("!level", user), "Level Command")
	assert.Equal(t, "", levelCommand.ProcessText("!level6", user), "Level no next command")
	levelCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", levelCommand.ProcessText("??", user), "Level next command")
}
