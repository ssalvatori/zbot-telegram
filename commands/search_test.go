package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

var searchCommand = SearchCommand{}

func TestSearchCommandOK(t *testing.T) {

	var result string

	searchCommand.Db = &db.ZbotDatabaseMock{
		SearchTerms: []string{"foo", "bar"},
	}

	result, _ = searchCommand.ProcessText("!search foo", userTest, "testchat", false)
	assert.Equal(t, "foo bar", result, "Search Command")

	searchCommand.Db = &db.ZbotDatabaseMock{
		SearchTerms: []string{},
	}

	result, _ = searchCommand.ProcessText("!search", userTest, "testchat", false)
	assert.Equal(t, "", result, "Search empty")
}
func TestSearchCommandNotMatch(t *testing.T) {

	result, err := searchCommand.ProcessText("!search6", userTest, "testchat", false)
	assert.Equal(t, "", result, "Empty output doesn't match")
	assert.Equal(t, err, ErrNextCommand, "Error output doesn't match")
}

func TestSearchCommandError(t *testing.T) {

	searchCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}
	_, err := searchCommand.ProcessText("!search foo", userTest, "testchat", false)
	assert.Error(t, err, "Internal Error")
}
