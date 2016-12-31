package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

var searchCommand = SearchCommand{}

func TestSearchCommandOK(t *testing.T) {

	searchCommand.Db = &db.MockZbotDatabase{
		Search_terms: []string{"foo", "bar",},
	}
	assert.Equal(t, "foo bar", searchCommand.ProcessText("!search foo", user), "Search Command")
	searchCommand.Db = &db.MockZbotDatabase{
		Search_terms: []string{},
	}
	assert.Equal(t, "", searchCommand.ProcessText("!search", user), "Search no next command")
	searchCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", searchCommand.ProcessText("?? ", user), "Search next command")
}