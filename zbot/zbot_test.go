package zbot

import (
	"testing"

	command "github.com/ssalvatori/zbot-telegram/commands"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestProcessingIsCommandDisabled(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	command.DisabledCommands = []string{
		"learn",
		"version",
	}

	botMsg := tb.Message{Text: "!learn", Sender: &tb.User{Username: "zbot_test"}}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "", result, "command disabled")

}

func TestProcessingVersion(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	buildTime = "2017-05-06 09:59:21.318841424 +0300 EEST"
	command.DisabledCommands = nil

	botMsg := tb.Message{
		Text: "!version",
		Sender: &tb.User{
			Username: "zbot_test",
		},
		Chat: &tb.Chat{
			Type:  "supergroup",
			Title: "testgroup",
		},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "zbot golang version ["+version+"] commit [undefined] build-time ["+buildTime+"]", result, "!version default")
}

func TestProcessingStats(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := tb.Message{
		Text: "!stats",
		Sender: &tb.User{
			Username: "zbot_test",
		},
		Chat: &tb.Chat{
			Type:  "supergroup",
			Title: "testgroup",
		},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, result, "Number of definitions: 666", "!stats")
}

func TestProcessingPing(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := tb.Message{
		Text: "!ping",
		Sender: &tb.User{
			Username: "zbot_test",
		},
		Chat: &tb.Chat{
			Type:  "supergroup",
			Title: "testgroup",
		},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, result, "pong!!", "!ping")
}

func TestProcessingRand(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "hola", Meaning: "gatolinux"}},
	}

	botMsg := tb.Message{Text: "!rand", Sender: &tb.User{Username: "zbot_test"}, Chat: &tb.Chat{Type: "private"}}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "[hola] - [gatolinux]", result, "!rand")
}

func TestProcessingGet(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level:   "666",
		File:    "hola.db",
		Term:    "hola",
		Meaning: "foo bar!",
	}

	botMsg := tb.Message{Text: "? hola", Sender: &tb.User{Username: "zbot_test"}, Chat: &tb.Chat{Type: "private"}}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, result, "[hola] - [foo bar!]", "? def fail")

}

func TestProcessingFind(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level:   "666",
		File:    "hola.db",
		Term:    "hola",
		Meaning: "foo bar!",
	}

	botMsg := tb.Message{Text: "!find hola", Sender: &tb.User{Username: "zbot_test"}, Chat: &tb.Chat{Type: "private"}}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, result, "hola", "!find fail")
}

func TestProcessingSearch(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level:       "666",
		File:        "hola.db",
		Term:        "hola",
		Meaning:     "foo bar!",
		FindTerms:   []string{"hola", "chao", "foo_bar"},
		RandDef:     []db.Definition{db.Definition{Term: "hola", Meaning: "gatolinux"}},
		SearchTerms: []string{"hola", "chao", "foobar"},
	}

	botMsg := tb.Message{Text: "!search hola", Sender: &tb.User{Username: "zbot_test"}, Chat: &tb.Chat{Type: "private"}}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "hola chao foobar", result, "!rand")
}

func TestProcessingUserLevel(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level:       "666",
		File:        "hola.db",
		Term:        "hola",
		Meaning:     "foo bar!",
		FindTerms:   []string{"hola", "chao", "foo_bar"},
		RandDef:     []db.Definition{db.Definition{Term: "hola", Meaning: "gatolinux"}},
		SearchTerms: []string{"hola", "chao", "foobar"},
	}

	botMsg := tb.Message{
		Text:   "!level",
		Sender: &tb.User{FirstName: "ssalvato", Username: "ssalvato"},
		Chat:   &tb.Chat{Type: "private"},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "ssalvato level 666", result, "!level self user")
}

func TestProcessingUserIgnoreList(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level:       "666",
		File:        "hola.db",
		Term:        "hola",
		Meaning:     "foo bar!",
		FindTerms:   []string{"hola", "chao", "foo_bar"},
		RandDef:     []db.Definition{db.Definition{Term: "hola", Meaning: "gatolinux"}},
		SearchTerms: []string{"hola", "chao", "foobar"},
		User_ignored: []db.UserIgnore{
			{Username: "ssalvato", CreatedAt: 1231, ValidUntil: 4564},
		},
	}

	botMsg := tb.Message{
		Text:   "!ignore list",
		Sender: &tb.User{FirstName: "ssalvato", Username: "ssalvato"},
		Chat:   &tb.Chat{Type: "private"},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "[ @ssalvato ] since [1970-01-01 00:20:31 +0000 UTC] until [1970-01-01 01:16:04 +0000 UTC]", result, "!ignore list")
}

func TestProcessingUserIgnoreInsert(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		FindTerms:    []string{"hola", "chao", "foo_bar"},
		RandDef:      []db.Definition{db.Definition{Term: "hola", Meaning: "gatolinux"}},
		SearchTerms:  []string{"hola", "chao", "foobar"},
		User_ignored: []db.UserIgnore{{Username: "ssalvatori", CreatedAt: 1231, ValidUntil: 4564}},
	}

	botMsg := tb.Message{
		Text:   "!ignore add rigo",
		Sender: &tb.User{FirstName: "ssalvatori", Username: "ssalvatori"},
		Chat:   &tb.Chat{Type: "private"},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "User [rigo] ignored for 10 minutes", result, "!ignore add OK")

	botMsg = tb.Message{
		Text:   "!ignore add ssalvatori",
		Sender: &tb.User{FirstName: "ssalvatori", Username: "ssalvatori"},
		Chat:   &tb.Chat{Type: "private"},
	}
	result = cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "You can't ignore yourself", result, "!ignore add myself")

}

func TestProcessingLearnReplyTo(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := tb.Message{Text: "!learn arg1",
		Sender: &tb.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tb.Message{
			Text: "message in reply-to",
			Sender: &tb.User{
				Username: "otheruser",
			},
		},
		Chat: &tb.Chat{Type: "private"},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)

	assert.Equal(t, "[arg1] - [otheruser message in reply-to]", result, "!learn with replayto")
}

func TestMessageProcessing(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	Flags.Ignore = false

	botMsg := tb.Message{Text: "!learn arg1",
		Sender: &tb.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tb.Message{
			Text: "message in reply-to",
			Sender: &tb.User{
				Username: "otheruser",
			},
		},
		Chat: &tb.Chat{Type: "private"},
	}

	result := cmdProcessing(dbMock, botMsg, "test_chat", false)

	assert.Equal(t, "[arg1] - [otheruser message in reply-to]", result, "!learn with replayto")
}

func TestMessagesProcessingIgnoredUser(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level:       "666",
		File:        "hola.db",
		Ignore_User: true,
	}

	Flags.Ignore = true

	botMsg := tb.Message{Text: "!learn arg1",
		Sender: &tb.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tb.Message{
			Text: "message in reply-to",
			Sender: &tb.User{
				Username: "otheruser",
			},
		},
		Chat: &tb.Chat{Type: "private"},
	}

	result := messagesProcessing(dbMock, &botMsg, "test_chat")
	assert.Equal(t, "", result, "!learn ignored")
}

func TestGetDisabledCommands(t *testing.T) {
	cmds := []string{"cmd1", "cmd2", "cmd3"}
	SetDisabledCommands(cmds)
	assert.Equal(t, cmds, GetDisabledCommands(), "Get Disabled Commands")

}

func TestProcessingNotEnoughPermissions(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level:       "666",
		File:        "hola.db",
		Ignore_User: true,
	}

	Flags.Level = true
	Flags.Ignore = false

	botMsg := tb.Message{Text: "!forget arg1",
		Sender: &tb.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tb.Message{
			Text: "message in reply-to",
			Sender: &tb.User{
				Username: "otheruser",
			},
		},
		Chat: &tb.Chat{Type: "private"},
	}

	result := messagesProcessing(dbMock, &botMsg, "test_chat")
	assert.Equal(t, "Your level is not enough < 1000", result, "Not enough permissions to use a command")
}

/*
func TestExecute(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level:             "666",
		File:              "hola.db",
		IgnoreListCleaned: false,
	}

	Flags.Ignore = true
	Execute()
	assert.Equal(t, true, dbMock.IgnoreListCleaned, "Ignore List Called")
}
*/
