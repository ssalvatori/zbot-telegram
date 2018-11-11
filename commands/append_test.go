package command

import (
	"fmt"
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var appendCommand = AppendCommand{}

func TestAppendCommandOK(t *testing.T) {

	appendCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level:   "100",
	}
	appendCommand.Levels = Levels{
		Ignore: 10,
		Append: 10,
		Learn:  10,
		Lock:   10,
	}

	userTest.Level = 100

	assert.Equal(t, "[foo] = [bar]", appendCommand.ProcessText("!append foo bar", userTest), "Append Command")

	appendCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level:   "5",
	}
	appendCommand.Levels = Levels{
		Ignore: 10,
		Append: 10,
		Learn:  10,
		Lock:   10,
	}

	userTest.Level = 5

	assert.Equal(t, fmt.Sprintf("Your level is not enough < %d", appendCommand.Levels.Lock), appendCommand.ProcessText("!append foo bar", userTest), "Append Command No Level")

	appendCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}
	appendCommand.Levels = Levels{
		Append: 1,
	}
	assert.Equal(t, "", appendCommand.ProcessText("!append foo bar2", userTest), "Append Error Set")

	appendCommand.Db = &db.MockZbotDatabase{
		ErrorAppend: true,
	}
	assert.Equal(t, "", appendCommand.ProcessText("!append foo bar2", userTest), "Append Error Get")

	appendCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", appendCommand.ProcessText("??", userTest), "Append next command")
}
