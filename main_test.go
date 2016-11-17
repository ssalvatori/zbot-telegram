package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/tucnak/telebot"
	"testing"
)

func TestProcessingPing(t *testing.T) {
	output := make(chan string)
	botMsg := telebot.Message{Text: "!ping"}
	go processing(botMsg, output)
	result := <-output
	assert.Equal(t, result, "pong!!", "!ping")
}

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

func TestProcessingVersion(t *testing.T) {
	output := make(chan string)
	botMsg := telebot.Message{Text: "!version sf"}
	go processing(botMsg, output)
	result := <-output
	assert.Equal(t, result, "zbot golang version 1.0", "!version fail")
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

/*
func TestMessagesProcessing(t *testing.T) {
	bot.Messages = make(chan telebot.Message)
	botMsg := [2]telebot.Message{ telebot.Message{Text: "!rand"}, telebot.Message{Text: "any text"}}
	bot.Messages <- botMsg
	go messagesProcessing()

}
*/
