package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var getCommand = GetCommand{}

func TestGetCommandOK(t *testing.T) {

	getCommand.Db = &db.MockZbotDatabase{
		Term: "foo",
		Meaning: "bar",
	}
	assert.Equal(t, "[foo] - [bar]", getCommand.ProcessText("? foo", user), "Last Command")
	getCommand.Db = &db.MockZbotDatabase{
		Not_found: true,
	}
	assert.Equal(t, "[foo2] Not found!", getCommand.ProcessText("? foo2", user), "Last no next command")
	getCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", getCommand.ProcessText("?? ", user), "Last next command")
}
