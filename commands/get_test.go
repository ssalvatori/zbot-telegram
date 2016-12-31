package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var command = GetCommand{}

func TestGetCommandOK(t *testing.T) {

	command.Db = &db.MockZbotDatabase{
		Term: "foo",
		Meaning: "bar",
	}
	assert.Equal(t, "[foo] - [bar]", command.ProcessText("? foo"), "Last Command")
	command.Db = &db.MockZbotDatabase{
		Not_found: true,
	}
	assert.Equal(t, "[foo2] Not found!", command.ProcessText("? foo2"), "Last no next command")
	command.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", command.ProcessText("?? "), "Last next command")
}
