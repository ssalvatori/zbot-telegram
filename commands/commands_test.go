package command

import (
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

var user = User{
	Username: "ssalvatori",
	Ident:    "stefano",
	Host:     "localhost",
}

type FakeCommand struct {
	Next HandlerCommand
}

func (handler *FakeCommand) ProcessText(text string, user User) string {
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
