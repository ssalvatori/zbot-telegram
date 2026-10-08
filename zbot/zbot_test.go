package zbot

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v4"
)

// testFlags is applied to every bot created with newTestBot. Tests mutate it before building.
var testFlags Flags

func newTestBot(d db.ZbotDatabase) *Bot {
	return newTestBotWithConfig(d, Config{})
}

func newTestBotWithConfig(d db.ZbotDatabase, cfg Config) *Bot {
	cfg.Flags = testFlags
	cfg.IgnoreDuration = 300
	cfg.Levels = DefaultLevels
	b, err := New(cfg, d)
	if err != nil {
		panic(err)
	}
	return b
}

func TestNew(t *testing.T) {
	_, err := New(Config{}, nil)
	assert.Error(t, err, "nil database is rejected")

	b, err := New(Config{}, &db.ZbotDatabaseMock{})
	assert.NoError(t, err)
	assert.NotEmpty(t, b.cfg.ModulesPath, "default modules path applied")
	assert.Equal(t, 300, b.cfg.IgnoreDuration, "default ignore duration applied")
	assert.NotNil(t, b.commands, "command chain built once at construction")
}

func TestProcessingIsCommandDisabled(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	b := newTestBotWithConfig(dbMock, Config{DisabledCommands: []string{"learn", "version"}})

	botMsg := tele.Message{Text: "!learn", Sender: &tele.User{Username: "zbot_test"}}
	result := b.cmdProcessing(botMsg, "test_chat", false)
	assert.Equal(t, "Command [learn] is disabled", result, "command disabled")

}

func TestProcessingVersion(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	buildTime = "2017-05-06 09:59:21.318841424 +0300 EEST"

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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
	assert.Equal(t, result, "pong!!", "!ping")
}

func TestProcessingRand(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		RandDef: []db.Definition{{Term: "hola", Meaning: "gatolinux"}},
	}

	botMsg := tele.Message{Text: "!rand", Sender: &tele.User{Username: "zbot_test"}, Chat: &tele.Chat{Type: "private"}}
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
	assert.Equal(t, "User [rigo] ignored for 5 minutes", result, "!ignore add OK")

	botMsg = tele.Message{
		Text:   "!ignore add ssalvatori",
		Sender: &tele.User{FirstName: "ssalvatori", Username: "ssalvatori"},
		Chat:   &tele.Chat{Type: "private"},
	}
	result = newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)

	assert.Equal(t, "[arg1] - [otheruser message in reply-to]", result, "!learn with replayto")
}

func TestMessageProcessing(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level: "666",
		File:  "hola.db",
	}

	testFlags.Ignore = false

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

	result := newTestBot(dbMock).cmdProcessing(botMsg, "test_chat", false)

	assert.Equal(t, "[arg1] - [otheruser message in reply-to]", result, "!learn with replayto")
}

func TestMessagesProcessingIgnoredUser(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level:      "666",
		File:       "hola.db",
		IgnoreUser: true,
	}

	testFlags.Ignore = true

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

	result := newTestBot(dbMock).messagesProcessing(&botMsg, "test_chat")
	assert.Equal(t, "", result, "!learn ignored")
}

func TestProcessingNotEnoughPermissions(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{
		Level:      "666",
		File:       "hola.db",
		IgnoreUser: true,
	}

	testFlags.Level = true
	testFlags.Ignore = false

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

	result := newTestBot(dbMock).messagesProcessing(&botMsg, "test_chat")
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

	b := newTestBot(&db.ZbotDatabaseMock{})

	msg := tele.Update{}
	assert.True(t, b.filterUpdate(&msg), "No Message")

	msg = tele.Update{Message: &tele.Message{Text: "test spam"}}
	assert.False(t, b.filterUpdate(&msg), "Spam is dropped")

	msg = tele.Update{Message: &tele.Message{Text: "test", Chat: &tele.Chat{Type: "private"}}}
	assert.True(t, b.filterUpdate(&msg), "Private message")
	assert.Empty(t, b.channels.Channels(), "private chats are not tracked")

	msg = tele.Update{Message: &tele.Message{Text: "test", Chat: &tele.Chat{Type: "group", ID: -1, Title: "g"}}}
	assert.True(t, b.filterUpdate(&msg), "Group message")
	assert.Equal(t, []Channel{{ID: -1, Title: "g"}}, b.channels.Channels(), "group chats are tracked")

}

func TestChannelRegistry(t *testing.T) {
	r := newChannelRegistry([]Channel{{Title: "preconfigured", AuthToken: "tok"}})

	r.Track(tele.Chat{ID: -1, Title: "preconfigured"})
	r.Track(tele.Chat{ID: -2, Title: "new"})

	got := r.Channels()
	assert.Equal(t, []Channel{{ID: -1, Title: "preconfigured", AuthToken: "tok"}, {ID: -2, Title: "new"}}, got)

	// The returned slice is a copy; mutating it must not affect the registry.
	got[0].Title = "changed"
	assert.Equal(t, "preconfigured", r.Channels()[0].Title)
}

func TestMessagesProcessingNoMatch(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	testFlags.Ignore = false

	// Message text doesn't start with ! or ? → empty result
	botMsg := tele.Message{
		Text:   "just a plain message",
		Sender: &tele.User{Username: "ssalvatori"},
		Chat:   &tele.Chat{Type: "group"},
	}
	result := newTestBot(dbMock).messagesProcessing(&botMsg, "testchat")
	assert.Equal(t, "", result)
}

func TestRunExternalModuleInvalidText(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	// Text that doesn't match /command pattern → ParseCommand returns error
	msg := &tele.Message{
		Text:   "not a slash command",
		Sender: &tele.User{Username: "user"},
		Chat:   &tele.Chat{Type: "group", Title: "test"},
	}
	result := newTestBot(dbMock).runExternalModule(msg)
	assert.Equal(t, "", result)
}

func TestRunExternalModuleCommandNotFound(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	// /cmd not present in empty modules list → GetCommandFile returns error
	msg := &tele.Message{
		Text:   "/unknowncmd arg",
		Sender: &tele.User{Username: "user"},
		Chat:   &tele.Chat{Type: "group", Title: "test"},
	}
	result := newTestBot(dbMock).runExternalModule(msg)
	assert.Equal(t, "", result)
}
