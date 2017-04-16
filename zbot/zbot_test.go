package zbot

import (
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
	"github.com/tucnak/telebot"

	"testing"
	"os"
)

func TestBuildUser(t *testing.T) {
	botMsg := telebot.Message{
		Sender: telebot.User{FirstName: "Stefano", Username: "Ssalvato"},
	}

	user := buildUser(botMsg.Sender)
	assert.Equal(t, "ssalvato", user.Username, "username defined")
	assert.Equal(t, "stefano", user.Ident, "ident defined")

	botMsg = telebot.Message{
		Sender: telebot.User{FirstName: "Stefano"},
	}

	user = buildUser(botMsg.Sender)
	assert.Equal(t, "stefano", user.Username, "username not defined")
	assert.Equal(t, user.Ident, "stefano", "ident defined")

}

func TestProcessingVersion(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := telebot.Message{Text: "!version"}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "zbot golang version ["+Version+"] build-time ["+BuildTime+"]", "!version default")
}

func TestProcessingStats(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := telebot.Message{Text: "!stats"}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "Count: 666", "!stats")
}

func TestProcessingPing(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := telebot.Message{Text: "!ping"}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "pong!!", "!ping")
}

func TestProcessingRand(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
	}

	botMsg := telebot.Message{Text: "!rand"}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "[hola] - [gatolinux]", result, "!rand")
}

func TestProcessingGet(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:   "666",
		File:    "hola.db",
		Term:    "hola",
		Meaning: "foo bar!",
	}

	botMsg := telebot.Message{Text: "? hola"}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "[hola] - [foo bar!]", "? def fail")

}

func TestProcessingFind(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:   "666",
		File:    "hola.db",
		Term:    "hola",
		Meaning: "foo bar!",
	}

	botMsg := telebot.Message{Text: "!find hola"}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "hola", "!find fail")
}

func TestProcessingSearch(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		Find_terms:   []string{"hola", "chao", "foo_bar"},
		Rand_def:     db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
		Search_terms: []string{"hola", "chao", "foobar"},
	}

	botMsg := telebot.Message{Text: "!search hola"}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "hola chao foobar", result, "!rand")
}

func TestProcessingUserLevel(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		Find_terms:   []string{"hola", "chao", "foo_bar"},
		Rand_def:     db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
		Search_terms: []string{"hola", "chao", "foobar"},
	}

	botMsg := telebot.Message{
		Text:   "!level",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "ssalvato level 666", result, "!rand")
}

func TestProcessingUserIgnoreList(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		Find_terms:   []string{"hola", "chao", "foo_bar"},
		Rand_def:     db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
		Search_terms: []string{"hola", "chao", "foobar"},
		User_ignored: []db.UserIgnore{
			{Username: "ssalvato", Since: "1231", Until: "4564"},
		},
	}

	botMsg := telebot.Message{
		Text:   "!ignore list",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "[ @ssalvato ] since [01-01-1970 00:20:31 UTC] until [01-01-1970 01:16:04 UTC]", result, "!ignore list")
}

func TestProcessingUserIgnoreInsert(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		Find_terms:   []string{"hola", "chao", "foo_bar"},
		Rand_def:     db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
		Search_terms: []string{"hola", "chao", "foobar"},
		User_ignored: []db.UserIgnore{{Username: "ssalvatori", Since: "1231", Until: "4564"}},
	}

	botMsg := telebot.Message{
		Text:   "!ignore add rigo",
		Sender: telebot.User{FirstName: "ssalvatori", Username: "ssalvatori"},
	}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "User [rigo] ignored for 10 minutes", result, "!ignore add OK")

	botMsg = telebot.Message{
		Text:   "!ignore add ssalvatori",
		Sender: telebot.User{FirstName: "ssalvatori", Username: "ssalvatori"},
	}
	result = processing(dbMock, botMsg)
	assert.Equal(t, "You can't ignore youself", result, "!ignore add myself")

	dbMock.Level = "10"
	botMsg = telebot.Message{
		Text:   "!ignore add ssalvato",
		Sender: telebot.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result = processing(dbMock, botMsg)
	assert.Equal(t, "Your level is not enough < 100", result, "!ignore")
}

func TestProcessingExternalModuleWithArgs(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := telebot.Message{Text: "!test arg1 arg2"}
	result := processing(dbMock, botMsg)

	assert.Equal(t, "OK arg1 arg2\n", result, "!test module with args")
}

func TestProcessingExternalModuleWithoutArgs(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := telebot.Message{Text: "!test"}
	result := processing(dbMock, botMsg)

	assert.Equal(t, "OK\n", result, "external module without args")
}

func TestProcessingExternalModuleNotFound(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := telebot.Message{Text: "!external arg1 arg2"}
	result := processing(dbMock, botMsg)

	assert.Equal(t, "", result, "external module not found")
}

func TestGetCurrentDirectory(t *testing.T) {
	directory := getCurrentDirectory()
	dir,_ := os.Getwd()
	assert.Equal(t, dir, directory,"getting current directory")
}

/*
func TestMessagesProcessing(t *testing.T) {
	dbMock := &db.MockZbotDatabase{
		Ignore_User: true,
	}
	msgChan := make(chan telebot.Message)
	bot := &telebot.Bot{Messages: msgChan}

	msgObj := telebot.Message{
		Text:   "!hola",
		Sender: telebot.User{FirstName: "Stefano", Username: "Ssalvato"},
	}
	bot.Messages <- msgObj
	go messagesProcessing(dbMock, bot)
}
*/
