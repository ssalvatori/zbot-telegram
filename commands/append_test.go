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
		Level: "100",
	}
	appendCommand.Levels = Levels{
		Ignore:10,
		Append:10,
		Learn:10,
		Lock:10,
	}

	assert.Equal(t, "[foo] = [bar]", appendCommand.ProcessText("!append foo bar", user), "Append Command")
}

func TestAppendCommandNoLevel(t *testing.T) {

	appendCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level: "5",
	}
	appendCommand.Levels = Levels{
		Ignore:10,
		Append:10,
		Learn:10,
		Lock:10,
	}

	assert.Equal(t, "", appendCommand.ProcessText("!append foo bar", user), "Append Command No Level")
}