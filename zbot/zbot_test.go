package zbot

import (
	"testing"

	command "github.com/ssalvatori/zbot-telegram/commands"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
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

	botMsg := tele.Message{Text: "!learn", Sender: &tele.User{Username: "zbot_test"}}
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

	botMsg := tele.Message{
		Text: "!version",
		Sender: &tele.User{
			Username: "zbot_test",
		},
		Chat: &tele.Chat{
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

	botMsg := tele.Message{
		Text: "!stats",
		Sender: &tele.User{
			Username: "zbot_test",
		},
		Chat: &tele.Chat{
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

	botMsg := tele.Message{
		Text: "!ping",
		Sender: &tele.User{
			Username: "zbot_test",
		},
		Chat: &tele.Chat{
			Type:  "supergroup",
			Title: "testgroup",
		},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, result, "pong!!", "!ping")
}

func TestProcessingRand(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		RandDef: []db.Definition{{Term: "hola", Meaning: "gatolinux"}},
	}

	botMsg := tele.Message{Text: "!rand", Sender: &tele.User{Username: "zbot_test"}, Chat: &tele.Chat{Type: "private"}}
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

	botMsg := tele.Message{Text: "? hola", Sender: &tele.User{Username: "zbot_test"}, Chat: &tele.Chat{Type: "private"}}
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

	botMsg := tele.Message{Text: "!find hola", Sender: &tele.User{Username: "zbot_test"}, Chat: &tele.Chat{Type: "private"}}
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
		RandDef:     []db.Definition{{Term: "hola", Meaning: "gatolinux"}},
		SearchTerms: []string{"hola", "chao", "foobar"},
	}

	botMsg := tele.Message{Text: "!search hola", Sender: &tele.User{Username: "zbot_test"}, Chat: &tele.Chat{Type: "private"}}
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
		RandDef:     []db.Definition{{Term: "hola", Meaning: "gatolinux"}},
		SearchTerms: []string{"hola", "chao", "foobar"},
	}

	botMsg := tele.Message{
		Text:   "!level",
		Sender: &tele.User{FirstName: "ssalvato", Username: "ssalvato"},
		Chat:   &tele.Chat{Type: "private"},
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
		RandDef:     []db.Definition{{Term: "hola", Meaning: "gatolinux"}},
		SearchTerms: []string{"hola", "chao", "foobar"},
		UserIgnored: []db.UserIgnore{
			{Username: "ssalvato", CreatedAt: 1231, ValidUntil: 4564},
		},
	}

	botMsg := tele.Message{
		Text:   "!ignore list",
		Sender: &tele.User{FirstName: "ssalvato", Username: "ssalvato"},
		Chat:   &tele.Chat{Type: "private"},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "[ @ssalvato ] since [1970-01-01 00:20:31 +0000 UTC] until [1970-01-01 01:16:04 +0000 UTC]", result, "!ignore list")
}

func TestProcessingUserIgnoreInsert(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level:       "666",
		File:        "hola.db",
		Term:        "hola",
		Meaning:     "foo bar!",
		FindTerms:   []string{"hola", "chao", "foo_bar"},
		RandDef:     []db.Definition{{Term: "hola", Meaning: "gatolinux"}},
		SearchTerms: []string{"hola", "chao", "foobar"},
		UserIgnored: []db.UserIgnore{{Username: "ssalvatori", CreatedAt: 1231, ValidUntil: 4564}},
	}

	botMsg := tele.Message{
		Text:   "!ignore add rigo",
		Sender: &tele.User{FirstName: "ssalvatori", Username: "ssalvatori"},
		Chat:   &tele.Chat{Type: "private"},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "User [rigo] ignored for 10 minutes", result, "!ignore add OK")

	botMsg = tele.Message{
		Text:   "!ignore add ssalvatori",
		Sender: &tele.User{FirstName: "ssalvatori", Username: "ssalvatori"},
		Chat:   &tele.Chat{Type: "private"},
	}
	result = cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "You can't ignore yourself", result, "!ignore add myself")

}

func TestProcessingLearnReplyTo(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := tele.Message{Text: "!learn arg1",
		Sender: &tele.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tele.Message{
			Text: "message in reply-to",
			Sender: &tele.User{
				Username: "otheruser",
			},
		},
		Chat: &tele.Chat{Type: "private"},
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

	botMsg := tele.Message{Text: "!learn arg1",
		Sender: &tele.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tele.Message{
			Text: "message in reply-to",
			Sender: &tele.User{
				Username: "otheruser",
			},
		},
		Chat: &tele.Chat{Type: "private"},
	}

	result := cmdProcessing(dbMock, botMsg, "test_chat", false)

	assert.Equal(t, "[arg1] - [otheruser message in reply-to]", result, "!learn with replayto")
}

func TestMessagesProcessingIgnoredUser(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level:      "666",
		File:       "hola.db",
		IgnoreUser: true,
	}

	Flags.Ignore = true

	botMsg := tele.Message{Text: "!learn arg1",
		Sender: &tele.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tele.Message{
			Text: "message in reply-to",
			Sender: &tele.User{
				Username: "otheruser",
			},
		},
		Chat: &tele.Chat{Type: "private"},
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
		Level:      "666",
		File:       "hola.db",
		IgnoreUser: true,
	}

	Flags.Level = true
	Flags.Ignore = false

	botMsg := tele.Message{Text: "!forget arg1",
		Sender: &tele.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tele.Message{
			Text: "message in reply-to",
			Sender: &tele.User{
				Username: "otheruser",
			},
		},
		Chat: &tele.Chat{Type: "private"},
	}

	result := messagesProcessing(dbMock, &botMsg, "test_chat")
	assert.Equal(t, "Your level is not enough < 1000", result, "Not enough permissions to use a command")
}

func TestAppendChannel(t *testing.T) {
	chat := &tele.Chat{
		Type:  "group",
		ID:    -1234,
		Title: "test 1",
	}

	channels := []Channel{}
	assert.Equal(t, []Channel{{ID: -1234, Title: "test 1"}}, appendChannel(channels, *chat), "Add Channel")

	channels = []Channel{{ID: -66, Title: "test 1"}}
	assert.Equal(t, []Channel{{ID: -66, Title: "test 1"}, {ID: -1234, Title: "test 1"}}, appendChannel(channels, *chat), "Add Channel")

	channels = []Channel{{ID: -1234, Title: "test already"}}
	assert.Equal(t, []Channel{{ID: -1234, Title: "test 1"}}, appendChannel(channels, *chat), "Channel already present (updating title)")

	channels = []Channel{{ID: -12345, Title: "test already"}, {ID: 0, Title: "test 1"}}
	assert.Equal(t, []Channel{{ID: -12345, Title: "test already"}, {ID: -1234, Title: "test 1"}}, appendChannel(channels, *chat), "Channel's ID is copied from message")
}

func TestMiddleware(t *testing.T) {

	msg := tele.Update{}
	assert.True(t, middlewareCustom(&msg), "No Message")

	msg = tele.Update{Message: &tele.Message{Text: "test spam"}}
	assert.False(t, middlewareCustom(&msg), "No Message")

	msg = tele.Update{Message: &tele.Message{Text: "test", Chat: &tele.Chat{Type: "private"}}}
	assert.True(t, middlewareCustom(&msg), "Private message")

	msg = tele.Update{Message: &tele.Message{Text: "test", Chat: &tele.Chat{Type: "group"}}}
	assert.True(t, middlewareCustom(&msg), "Group message")

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
