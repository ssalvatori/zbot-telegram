package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var searchCommand = SearchCommand{}

func TestSearchCommandOK(t *testing.T) {

	searchCommand.Db = &db.MockZbotDatabase{
		Search_terms: []string{"foo", "bar"},
	}
	assert.Equal(t, "foo bar", searchCommand.ProcessText("!search foo", userTest), "Search Command")
	searchCommand.Db = &db.MockZbotDatabase{
		Search_terms: []string{},
	}
	assert.Equal(t, "", searchCommand.ProcessText("!search", userTest), "Search no next command")
	searchCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", searchCommand.ProcessText("?? ", userTest), "Search next command")
}
