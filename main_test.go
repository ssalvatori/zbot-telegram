package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/tucnak/telebot"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

func TestProcessingVersion(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!version"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, result, "zbot golang version 1.0", "!version fail")
}

func TestProcessingStats(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!stats"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, result, "Count: 666", "!stats")
}

func TestProcessingPing(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!ping"}
	result := processing(dbMock,botMsg, output)
	assert.Equal(t, result, "pong!!", "!ping")
}

func TestProcessingRand(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "hola", Meaning:"gatolinux"},
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!rand"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "[hola] - [gatolinux]", result,  "!rand")
}

func TestProcessingGet(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "? hola"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, result, "[hola] - [foo bar!]", "? def fail")

}

func TestProcessingFind(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!find hola"}
	result := processing(dbMock ,botMsg, output)
	assert.Equal(t, result, "hola", "!find fail")
}

func TestProcessingSearch(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
		Find_terms: []string{"hola", "chao", "foo_bar",},
		Rand_def: db.DefinitionItem{Term: "hola", Meaning:"gatolinux"},
		Search_terms: []string{"hola","chao", "foobar"},
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!search hola"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "hola chao foobar", result,  "!rand")
}
/*
func TestGetTerms(t *testing.T) {

	var items = []db.DefinitionItem {
		{Term: "Term1"},
		{Term: "Term2"},
		{Meaning: ""},
	}

	terms := getTerms(items)
	assert.Equal(t, terms, []string{"Term1", "Term2"} )
}

func TestGetUserIgnored(t *testing.T) {
	var users = []db.UserIgnore {
		{
			Username: "rigo",
			Since: "1478126960",
			Until: "1478127560",
		},
	}

	formated := getUserIgnored(users)
	assert.Equal(t, formated, []string{"[ @rigo ] since [1478126960] until [1478127560]"})
}













func TestProcessingTop(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
		Find_terms: []string{"hola", "chao", "foo_bar",},
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!top"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, result, "hola chao foo_bar", "!top")
}



func TestProcessingLearn(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
		Find_terms: []string{"hola", "chao", "foo_bar",},
		Rand_def: db.DefinitionItem{Term: "hola", Meaning:"gatolinux"},
	}

	output := make(chan string)
	botMsg := telebot.Message{
		Text: "!learn 12312 foo bar!",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "[12312] - [foo bar!]", result, "!learn fail")
}




func TestProcessingLast(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
		Find_terms: []string{"hola", "chao", "foo_bar",},
		Rand_def: db.DefinitionItem{Term: "hola", Meaning:"gatolinux"},
		Search_terms: []string{"hola","chao", "foobar"},
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!last"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "[hola] - [foo bar!]", result,  "!rand")
}

func TestProcessingUserLevel(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
		Find_terms: []string{"hola", "chao", "foo_bar",},
		Rand_def: db.DefinitionItem{Term: "hola", Meaning:"gatolinux"},
		Search_terms: []string{"hola","chao", "foobar"},
	}

	output := make(chan string)
	botMsg := telebot.Message{
		Text: "!level",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "ssalvato level 666", result,  "!rand")
}

func TestProcessingUserIgnoreList(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
		Find_terms: []string{"hola", "chao", "foo_bar",},
		Rand_def: db.DefinitionItem{Term: "hola", Meaning:"gatolinux"},
		Search_terms: []string{"hola","chao", "foobar"},
		User_ignored: []db.UserIgnore{db.UserIgnore{Username: "ssalvato", Since:"1231", Until: "4564"},},
	}


	output := make(chan string)
	botMsg := telebot.Message{
		Text: "!ignorelist",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "[ @ssalvato ] since [1231] until [4564]", result,  "!rand")
}

func TestProcessingUserIgnoreInsert(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File: "hola.db",
		Term: "hola",
		Meaning: "foo bar!",
		Find_terms: []string{"hola", "chao", "foo_bar",},
		Rand_def: db.DefinitionItem{Term: "hola", Meaning:"gatolinux"},
		Search_terms: []string{"hola","chao", "foobar"},
		User_ignored: []db.UserIgnore{db.UserIgnore{Username: "ssalvato", Since:"1231", Until: "4564"},},
	}


	output := make(chan string)
	botMsg := telebot.Message{
		Text: "!ignore rigo",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "User [rigo] ignored for 10 minutes", result,  "!rand")

	output = make(chan string)
	botMsg = telebot.Message{
		Text: "!ignore ssalvato",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result = processing(dbMock, botMsg, output)
	assert.Equal(t, "You can't ignore youself", result,  "!ignore")


	dbMock.Level = "10"
	output = make(chan string)
	botMsg = telebot.Message{
		Text: "!ignore ssalvato",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result = processing(dbMock, botMsg, output)
	assert.Equal(t, "level not enough (minimum 500 yours 10)", result,  "!ignore")
}*/
