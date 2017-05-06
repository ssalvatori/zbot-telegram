package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"github.com/stretchr/testify/assert"
)

// This userTest variable will be shared between all the test
var userTest = user.User{
	Username: "ssalvatori",
	Ident:    "stefano",
	Host:     "localhost",
	Level:    100,
}

type FakeCommand struct {
	Next HandlerCommand
}

func (handler *FakeCommand) ProcessText(text string, user user.User) string {
	return "Fake OK"
}

func TestGetTerms(t *testing.T) {
	items := []db.DefinitionItem{
		{Term: "foo", Meaning: "bar"},
		{Term: "foo2", Meaning: "bar2"},
		{Term: "foo3", Meaning: "bar3"},
	}
	assert.Equal(t, []string{"foo", "foo2", "foo3"}, getTerms(items))
	var terms []string
	assert.Equal(t, terms, getTerms([]db.DefinitionItem{}))
}
