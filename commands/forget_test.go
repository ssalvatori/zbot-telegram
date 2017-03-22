package command

import (
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var forgetCommand = ForgetCommand{}

func TestForgetCommandOK(t *testing.T) {

	forgetCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level: "100",
	}
	forgetCommand.Levels = Levels{
		Ignore:10,
		Append:10,
		Learn:10,
		Lock:10,
		Forget: 10,
	}

	assert.Equal(t, "[foo] deleted", forgetCommand.ProcessText("!forget foo", user), "Forget Command OK")
}

func TestForgetCommandNoLevel(t *testing.T) {

	forgetCommand.Db = &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
		Level: "5",
	}
	forgetCommand.Levels = Levels{
		Ignore:10,
		Append:10,
		Learn:10,
		Lock:10,
		Forget: 1000,
	}

	assert.Equal(t, "", forgetCommand.ProcessText("!forget foo", user), "Forget Command No Level")
}