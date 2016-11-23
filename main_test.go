package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/tucnak/telebot"
)

func TestGetTerms(t *testing.T) {

	var items = []definitionItem {
		{term: "Term1"},
		{term: "Term2"},
		{meaning: ""},
	}

	terms := getTerms(items)
	assert.Equal(t, terms, []string{"Term1", "Term2"} )
}

func TestGetUserIgnored(t *testing.T) {
	var users = []userIgnore {
		{
			username: "rigo",
			since: "1478126960",
			until: "1478127560",
		},
	}

	formated := getUserIgnored(users)
	assert.Equal(t, formated, []string{"[ @rigo ] since [1478126960] until [1478127560]"})
}

func TestProcessingVersion(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!version"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, result, "zbot golang version 1.0", "!version fail")
}


func TestProcessingPing(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!ping"}
	result := processing(dbMock,botMsg, output)
	assert.Equal(t, result, "pong!!", "!ping")
}


func TestProcessingStats(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!stats"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, result, "Count: 666", "!stats")
}

func TestProcessingGet(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "? hola"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, result, "[hola] - [foo bar!]", "? def fail")

}

func TestProcessingFind(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!find hola"}
	result := processing(dbMock ,botMsg, output)
	assert.Equal(t, result, "hola", "!find fail")
}

func TestProcessingTop(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
		find_terms: []string{"hola", "chao", "foo_bar",},
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!top"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, result, "hola chao foo_bar", "!top")
}

func TestProcessingRand(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
		find_terms: []string{"hola", "chao", "foo_bar",},
		rand_def: definitionItem{term: "hola", meaning:"gatolinux"},
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!rand"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "[hola] - [gatolinux]", result,  "!rand")
}

func TestProcessingLearn(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
		find_terms: []string{"hola", "chao", "foo_bar",},
		rand_def: definitionItem{term: "hola", meaning:"gatolinux"},
	}

	output := make(chan string)
	botMsg := telebot.Message{
		Text: "!learn 12312 foo bar!",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "[12312] - [foo bar!]", result, "!learn fail")
}


func TestProcessingSearch(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
		find_terms: []string{"hola", "chao", "foo_bar",},
		rand_def: definitionItem{term: "hola", meaning:"gatolinux"},
		search_terms: []string{"hola","chao", "foobar"},
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!search hola"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "hola chao foobar", result,  "!rand")
}

func TestProcessingLast(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
		find_terms: []string{"hola", "chao", "foo_bar",},
		rand_def: definitionItem{term: "hola", meaning:"gatolinux"},
		search_terms: []string{"hola","chao", "foobar"},
	}

	output := make(chan string)
	botMsg := telebot.Message{Text: "!last"}
	result := processing(dbMock, botMsg, output)
	assert.Equal(t, "[hola] - [foo bar!]", result,  "!rand")
}

func TestProcessingUserLevel(t *testing.T) {

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
		find_terms: []string{"hola", "chao", "foo_bar",},
		rand_def: definitionItem{term: "hola", meaning:"gatolinux"},
		search_terms: []string{"hola","chao", "foobar"},
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

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
		find_terms: []string{"hola", "chao", "foo_bar",},
		rand_def: definitionItem{term: "hola", meaning:"gatolinux"},
		search_terms: []string{"hola","chao", "foobar"},
		user_ignored: []userIgnore{userIgnore{username: "ssalvato", since:"1231", until: "4564"},},
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

	dbMock := &mockZbotDatabase{
		level: "666",
		file: "hola.db",
		term: "hola",
		meaning: "foo bar!",
		find_terms: []string{"hola", "chao", "foo_bar",},
		rand_def: definitionItem{term: "hola", meaning:"gatolinux"},
		search_terms: []string{"hola","chao", "foobar"},
		user_ignored: []userIgnore{userIgnore{username: "ssalvato", since:"1231", until: "4564"},},
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


	dbMock.level = "10"
	output = make(chan string)
	botMsg = telebot.Message{
		Text: "!ignore ssalvato",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result = processing(dbMock, botMsg, output)
	assert.Equal(t, "level not enough (minimum 500 yours 10)", result,  "!ignore")
}