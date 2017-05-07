package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var getCommand = GetCommand{}

func TestGetCommandOK(t *testing.T) {

	getCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
	}

	assert.Equal(t, "[foo] - [bar]", getCommand.ProcessText("? foo", userTest), "Last Command")

	getCommand.Db = &db.MockZbotDatabase{
		Not_found: true,
	}

	assert.Equal(t, "[foo2] Not found!", getCommand.ProcessText("? foo2", userTest), "Last no next command")

	getCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", getCommand.ProcessText("?? ", userTest), "Last next command")
}
