package command

import (
"testing"
"github.com/ssalvatori/zbot-telegram-go/db"
"github.com/stretchr/testify/assert"
	"fmt"
)

var lockCommand = LockCommand{}

func TestTemplateCommandOK(t *testing.T) {

	lockCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level: "100",
	}
	lockCommand.Levels = Levels{
		Ignore:10,
		Append:10,
		Learn:10,
		Lock:10,
	}

	assert.Equal(t, "[foo] locked", lockCommand.ProcessText("!lock foo", user), "Template Command")
}

func TestTemplateCommandNoLevel(t *testing.T) {

	lockCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level: "5",
	}
	lockCommand.Levels = Levels{
		Ignore:10,
		Append:10,
		Learn:10,
		Lock:100,
	}

	assert.Equal(t, fmt.Sprintf("Your level is not enough < %s", lockCommand.Levels.Lock), lockCommand.ProcessText("!lock foo", user), "Lock Command No Level")
}