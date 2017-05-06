package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var levelCommand = LevelCommand{}

func TestLevelCommandOK(t *testing.T) {

	levelCommand.Db = &db.MockZbotDatabase{
		Level: "1000",
	}

	assert.Equal(t, "ssalvatori level 1000", levelCommand.ProcessText("!level", userTest), "Get Level from the same user")
}

func TestLevelAdd(t *testing.T) {

	levelCommand.Db = &db.MockZbotDatabase{
		Level: "10",
	}

	levelCommand.Levels = Levels{
		LevelAdd: 100,
	}

	assert.Equal(t, "not ready", levelCommand.ProcessText("!level add rigo 10", userTest), "add user")

}

func TestLevelOthers(t *testing.T) {

	levelCommand.Db = &db.MockZbotDatabase{
		Level: "10",
	}

	assert.Equal(t, "", levelCommand.ProcessText("!level6", userTest), "Level no next command")
	levelCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", levelCommand.ProcessText("??", userTest), "Level next command")
}
