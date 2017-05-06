package command

import (
	"fmt"
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var lockCommand = LockCommand{}

func TestTemplateCommandOK(t *testing.T) {

	lockCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level:   "100",
	}
	lockCommand.Levels = Levels{
		Ignore: 10,
		Append: 10,
		Learn:  10,
		Lock:   1,
	}

	assert.Equal(t, "[foo] locked", lockCommand.ProcessText("!lock foo", userTest), "Template Command")
}

func TestTemplateCommandNoLevel(t *testing.T) {

	lockCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level:   "5",
	}
	lockCommand.Levels = Levels{
		Ignore: 10,
		Append: 10,
		Learn:  10,
		Lock:   100,
	}

	userTest.Level = 5

	assert.Equal(t, fmt.Sprintf("Your level is not enough < %s", lockCommand.Levels.Lock), lockCommand.ProcessText("!lock foo", userTest), "Lock Command No Level")
}
