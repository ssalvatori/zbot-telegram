package command

import (
"testing"
"github.com/ssalvatori/zbot-telegram-go/db"
"github.com/stretchr/testify/assert"
)

var whoCommand = WhoCommand{}

func TestWhoCommandOK(t *testing.T) {

	whoCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Author: "ssalvatori",
		Date: "2017-03-22",
		Level: "100",
	}
	whoCommand.Levels = Levels{
		Ignore:10,
		Append:10,
		Learn:10,
		Lock:10,
	}

	assert.Equal(t, "[foo] by [ssalvatori] on [2017-03-22]", whoCommand.ProcessText("!who foo", user), "Who Command OK")
}

func TestWhoCommandNoLevel(t *testing.T) {

	whoCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level: "5",
	}
	whoCommand.Levels = Levels{
		Ignore:10,
		Append:10,
		Learn:10,
		Lock:10,
	}

	assert.Equal(t, "", appendCommand.ProcessText("!who foo", user), "Who Command No Level")
}


