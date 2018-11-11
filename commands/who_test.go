package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var whoCommand = WhoCommand{}

func TestWhoCommand(t *testing.T) {

	whoCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Author:  "ssalvatori",
		Date:    "2017-03-22",
		Level:   "100",
	}
	whoCommand.Levels = Levels{
		Who: 1,
	}

	assert.Equal(t, "[foo] by [ssalvatori] on [2017-03-22]", whoCommand.ProcessText("!who foo", userTest), "Who Command OK")

	whoCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level:   "5",
	}
	whoCommand.Levels = Levels{
		Who: 10,
	}

	assert.Equal(t, "", whoCommand.ProcessText("!who foo", userTest), "Who Command No Level")

	whoCommand.Db = &db.MockZbotDatabase{
		Error: true,
	}
	whoCommand.Levels = Levels{
		Who: 1,
	}
	assert.Equal(t, "", whoCommand.ProcessText("!who foo", userTest), "Who with internal error")

	whoCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", whoCommand.ProcessText("!who2 foo", userTest), "Who next command")
}
