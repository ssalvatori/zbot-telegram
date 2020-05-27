package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var searchCommand = SearchCommand{}

func TestSearchCommandOK(t *testing.T) {

	var result string

	searchCommand.Db = &db.MockZbotDatabase{
		Search_terms: []string{"foo", "bar"},
	}

	result, _ = searchCommand.ProcessText("!search foo", userTest)
	assert.Equal(t, "foo bar", result, "Search Command")

	searchCommand.Db = &db.MockZbotDatabase{
		Search_terms: []string{},
	}

	result, _ = searchCommand.ProcessText("!search", userTest)
	assert.Equal(t, "", result, "Search empty")
}
func TestSearchCommandNotMatch(t *testing.T) {

	result, _ := searchCommand.ProcessText("!search6", userTest)
	assert.Equal(t, "", result, "Empty output doesn't match")

	_, err := searchCommand.ProcessText("!search6", userTest)
	assert.Equal(t, "text doesn't match", err.Error(), "Error output doesn't match")
}

func TestSearchCommandError(t *testing.T) {

	searchCommand.Db = &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
		Error:    true,
	}
	_, err := searchCommand.ProcessText("!search foo", userTest)
	assert.Equal(t, "mock", err.Error(), "Db error")
}
