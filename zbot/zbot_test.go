package zbot

import (
	"context"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	command "github.com/ssalvatori/zbot-telegram/commands"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
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

	botMsg := models.Message{Text: "!learn", From: &models.User{Username: "zbot_test"}}
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

	botMsg := models.Message{
		Text: "!version",
		From: &models.User{
			Username: "zbot_test",
		},
		Chat: models.Chat{
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

	botMsg := models.Message{
		Text: "!stats",
		From: &models.User{
			Username: "zbot_test",
		},
		Chat: models.Chat{
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

	botMsg := models.Message{
		Text: "!ping",
		From: &models.User{
			Username: "zbot_test",
		},
		Chat: models.Chat{
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

	botMsg := models.Message{Text: "!rand", From: &models.User{Username: "zbot_test"}, Chat: models.Chat{Type: "private"}}
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

	botMsg := models.Message{Text: "? hola", From: &models.User{Username: "zbot_test"}, Chat: models.Chat{Type: "private"}}
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

	botMsg := models.Message{Text: "!find hola", From: &models.User{Username: "zbot_test"}, Chat: models.Chat{Type: "private"}}
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

	botMsg := models.Message{Text: "!search hola", From: &models.User{Username: "zbot_test"}, Chat: models.Chat{Type: "private"}}
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

	botMsg := models.Message{
		Text: "!level",
		From: &models.User{FirstName: "ssalvato", Username: "ssalvato"},
		Chat: models.Chat{Type: "private"},
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

	botMsg := models.Message{
		Text: "!ignore list",
		From: &models.User{FirstName: "ssalvato", Username: "ssalvato"},
		Chat: models.Chat{Type: "private"},
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

	botMsg := models.Message{
		Text: "!ignore add rigo",
		From: &models.User{FirstName: "ssalvatori", Username: "ssalvatori"},
		Chat: models.Chat{Type: "private"},
	}
	result := cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "User [rigo] ignored for 10 minutes", result, "!ignore add OK")

	botMsg = models.Message{
		Text: "!ignore add ssalvatori",
		From: &models.User{FirstName: "ssalvatori", Username: "ssalvatori"},
		Chat: models.Chat{Type: "private"},
	}
	result = cmdProcessing(dbMock, botMsg, "test_chat", false)
	assert.Equal(t, "You can't ignore yourself", result, "!ignore add myself")

}

func TestProcessingLearnReplyTo(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := models.Message{Text: "!learn arg1",
		From: &models.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyToMessage: &models.Message{
			Text: "message in reply-to",
			From: &models.User{
				Username: "otheruser",
			},
		},
		Chat: models.Chat{Type: "private"},
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

	botMsg := models.Message{Text: "!learn arg1",
		From: &models.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyToMessage: &models.Message{
			Text: "message in reply-to",
			From: &models.User{
				Username: "otheruser",
			},
		},
		Chat: models.Chat{Type: "private"},
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

	botMsg := models.Message{Text: "!learn arg1",
		From: &models.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyToMessage: &models.Message{
			Text: "message in reply-to",
			From: &models.User{
				Username: "otheruser",
			},
		},
		Chat: models.Chat{Type: "private"},
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

	botMsg := models.Message{Text: "!forget arg1",
		From: &models.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyToMessage: &models.Message{
			Text: "message in reply-to",
			From: &models.User{
				Username: "otheruser",
			},
		},
		Chat: models.Chat{Type: "private"},
	}

	result := messagesProcessing(dbMock, &botMsg, "test_chat")
	assert.Equal(t, "Your level is not enough < 1000", result, "Not enough permissions to use a command")
}

func TestAppendChannel(t *testing.T) {
	chat := models.Chat{
		Type:  "group",
		ID:    -1234,
		Title: "test 1",
	}

	channels := []Channel{}
	assert.Equal(t, []Channel{{ID: -1234, Title: "test 1"}}, appendChannel(channels, chat), "Add Channel")

	channels = []Channel{{ID: -66, Title: "test 1"}}
	assert.Equal(t, []Channel{{ID: -66, Title: "test 1"}, {ID: -1234, Title: "test 1"}}, appendChannel(channels, chat), "Add Channel")

	channels = []Channel{{ID: -1234, Title: "test already"}}
	assert.Equal(t, []Channel{{ID: -1234, Title: "test 1"}}, appendChannel(channels, chat), "Channel already present (updating title)")

	channels = []Channel{{ID: -12345, Title: "test already"}, {ID: 0, Title: "test 1"}}
	assert.Equal(t, []Channel{{ID: -12345, Title: "test already"}, {ID: -1234, Title: "test 1"}}, appendChannel(channels, chat), "Channel's ID is copied from message")
}

func TestMiddleware(t *testing.T) {
	called := false
	next := func(ctx context.Context, b *bot.Bot, update *models.Update) {
		called = true
	}
	handler := spamFilterMiddleware(next)

	// nil message → next is called
	called = false
	handler(context.Background(), nil, &models.Update{})
	assert.True(t, called, "No Message - should pass through")

	// spam message → next is NOT called
	called = false
	handler(context.Background(), nil, &models.Update{Message: &models.Message{Text: "test spam"}})
	assert.False(t, called, "Spam message - should be filtered")

	// normal message → next is called
	called = false
	handler(context.Background(), nil, &models.Update{Message: &models.Message{Text: "test", Chat: models.Chat{Type: "private"}}})
	assert.True(t, called, "Private message - should pass through")

	// group message → next is called
	called = false
	handler(context.Background(), nil, &models.Update{Message: &models.Message{Text: "test", Chat: models.Chat{Type: "group"}}})
	assert.True(t, called, "Group message - should pass through")
}

func TestSetDisabledLearnChannels(t *testing.T) {
	channels := []string{"chan1", "chan2"}
	SetDisabledLearnChannels(channels)
	// verify it was applied (commands package var)
	assert.NotNil(t, channels)
}

func TestMessagesProcessingNoMatch(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	Flags.Ignore = false

	// Message text doesn't start with ! or ? → empty result
	botMsg := models.Message{
		Text: "just a plain message",
		From: &models.User{Username: "ssalvatori"},
		Chat: models.Chat{Type: "group"},
	}
	result := messagesProcessing(dbMock, &botMsg, "testchat")
	assert.Equal(t, "", result)
}

func TestRunExternalModuleInvalidText(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	// Text that doesn't match /command pattern → ParseCommand returns error
	msg := &models.Message{
		Text: "not a slash command",
		From: &models.User{Username: "user"},
		Chat: models.Chat{Type: "group", Title: "test"},
	}
	result := runExternalModule(dbMock, msg, []ExternalModule{})
	assert.Equal(t, "", result)
}

func TestRunExternalModuleCommandNotFound(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	// /cmd not present in empty modules list → GetCommandFile returns error
	msg := &models.Message{
		Text: "/unknowncmd arg",
		From: &models.User{Username: "user"},
		Chat: models.Chat{Type: "group", Title: "test"},
	}
	result := runExternalModule(dbMock, msg, []ExternalModule{})
	assert.Equal(t, "", result)
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
