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


/*
func TestProcessingLearn(t *testing.T) {

	output := make(chan string)
	botMsg := telebot.Message{Text: "!learn 12312 foo bar!"}
	go processing(botMsg, output)
	result := <-output
	assert.Equal(t, result, "[12312] - [foo bar!]", "!learn fail")
}

func TestProcessingGet(t *testing.T) {

	output := make(chan string)
	botMsg := telebot.Message{Text: "? hola"}
	go processing(botMsg, output)
	result := <-output
	assert.Equal(t, result, "[hola] - [foo bar!]", "? def fail")

}



func TestProcessingFind(t *testing.T) {
	output := make(chan string)
	botMsg := telebot.Message{Text: "!find hola"}
	go processing(botMsg, output)
	result := <-output
	assert.Equal(t, result, "hola", "!find fail")
}

func TestProcessingTop(t *testing.T) {
	output := make(chan string)
	botMsg := telebot.Message{Text: "!top"}
	go processing(botMsg, output)
	result := <-output
	assert.Equal(t, result, "hola chao foo_bar", "!top")
}

func TestProcessingRand(t *testing.T) {
	output := make(chan string)
	botMsg := telebot.Message{Text: "!rand"}
	go processing(botMsg, output)
	result := <-output
	assert.Equal(t, result, "[hola] - [gatolinux]", "!rand")
}

func TestProcessingStats(t *testing.T) {
	output := make(chan string)
	botMsg := telebot.Message{Text: "!stats"}
	go processing(botMsg, output)
	result := <-output
	assert.Equal(t, result, "Count: 17461", "!stats")
}

func TestMessagesProcessing(t *testing.T) {
	bot.Messages = make(chan telebot.Message)
	botMsg := [2]telebot.Message{ telebot.Message{Text: "!rand"}, telebot.Message{Text: "any text"}}
	bot.Messages <- botMsg
	go messagesProcessing()

}
*/
