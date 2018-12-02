package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var appendCommand = AppendCommand{}

func TestAppendCommandOK(t *testing.T) {

	appendCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
	}

	assert.Equal(t, "[foo] = [bar]", appendCommand.ProcessText("!append foo bar", userTest), "Append Command")

	appendCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}

	assert.Equal(t, "", appendCommand.ProcessText("!append foo bar2", userTest), "Append Error Set")

	appendCommand.Db = &db.MockZbotDatabase{
		ErrorAppend: true,
	}
	assert.Equal(t, "", appendCommand.ProcessText("!append foo bar2", userTest), "Append Error Get")

	appendCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", appendCommand.ProcessText("??", userTest), "Append next command")
}
