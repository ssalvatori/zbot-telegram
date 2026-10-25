package zbot

import (
	"testing"
	"time"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v4"
)

// testFlags is applied to every bot created with newTestBot. Tests mutate it before building.
var testFlags Flags

// testChannel is the resolved channel handed to cmdProcessing in tests.
var testChannel = db.Channel{ID: 1, TelegramID: -1, Name: "test_chat", LearnEnabled: true, ExternalEnabled: true}

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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
	assert.Equal(t, "Number of definitions: 666", result, "!stats")
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
	assert.Equal(t, result, "pong!!", "!ping")
}

func TestProcessingRand(t *testing.T) {

	dbMock := &db.ZbotDatabaseMock{
		RandDef: []db.Definition{{Term: "hola", Meaning: "gatolinux"}},
	}

	botMsg := tele.Message{Text: "!rand", Sender: &tele.User{Username: "zbot_test"}, Chat: &tele.Chat{Type: "private"}}
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
	assert.Equal(t, "User [rigo] ignored for 5 minutes", result, "!ignore add OK")

	botMsg = tele.Message{
		Text:   "!ignore add ssalvatori",
		Sender: &tele.User{FirstName: "ssalvatori", Username: "ssalvatori"},
		Chat:   &tele.Chat{Type: "private"},
	}
	result = newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)
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
	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)

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

	result := newTestBot(dbMock).cmdProcessing(botMsg, testChannel, false)

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
		Chat: &tele.Chat{Type: "group", ID: -1, Title: "test_chat"},
	}

	result := newTestBot(dbMock).messagesProcessing(&botMsg)
	assert.Equal(t, "", result, "!learn ignored")
}

func TestMessagesProcessingChannelNotFound(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{Error: true}
	testFlags.Ignore = false

	botMsg := tele.Message{
		Text:   "!ping",
		Sender: &tele.User{Username: "ssalvatori"},
		Chat:   &tele.Chat{Type: "group", ID: -1, Title: "unknown"},
	}
	result := newTestBot(dbMock).messagesProcessing(&botMsg)
	assert.Equal(t, "", result, "unknown channel is dropped")
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
		Chat: &tele.Chat{Type: "group", ID: -1, Title: "test_chat"},
	}

	result := newTestBot(dbMock).messagesProcessing(&botMsg)
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
		Chat:   &tele.Chat{Type: "group", ID: -1, Title: "testchat"},
	}
	result := newTestBot(dbMock).messagesProcessing(&botMsg)
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

// addedToGroupCtx implements tele.Context for testing handleAddedToGroup.
type addedToGroupCtx struct {
	sender *tele.User
	chat   *tele.Chat
	msg    *tele.Message
	sent   string
}

func (m *addedToGroupCtx) Bot() tele.API                                { return nil }
func (m *addedToGroupCtx) Update() tele.Update                          { return tele.Update{} }
func (m *addedToGroupCtx) Message() *tele.Message                       { return m.msg }
func (m *addedToGroupCtx) Callback() *tele.Callback                     { return nil }
func (m *addedToGroupCtx) Query() *tele.Query                           { return nil }
func (m *addedToGroupCtx) InlineResult() *tele.InlineResult             { return nil }
func (m *addedToGroupCtx) ShippingQuery() *tele.ShippingQuery           { return nil }
func (m *addedToGroupCtx) PreCheckoutQuery() *tele.PreCheckoutQuery     { return nil }
func (m *addedToGroupCtx) Payment() *tele.Payment                       { return nil }
func (m *addedToGroupCtx) Poll() *tele.Poll                             { return nil }
func (m *addedToGroupCtx) PollAnswer() *tele.PollAnswer                 { return nil }
func (m *addedToGroupCtx) ChatMember() *tele.ChatMemberUpdate           { return nil }
func (m *addedToGroupCtx) ChatJoinRequest() *tele.ChatJoinRequest       { return nil }
func (m *addedToGroupCtx) Migration() (int64, int64)                    { return 0, 0 }
func (m *addedToGroupCtx) Topic() *tele.Topic                           { return nil }
func (m *addedToGroupCtx) Boost() *tele.BoostUpdated                    { return nil }
func (m *addedToGroupCtx) BoostRemoved() *tele.BoostRemoved             { return nil }
func (m *addedToGroupCtx) PurchasedPaidMedia() *tele.PaidMediaPurchased { return nil }
func (m *addedToGroupCtx) Sender() *tele.User                           { return m.sender }
func (m *addedToGroupCtx) Chat() *tele.Chat                             { return m.chat }
func (m *addedToGroupCtx) Recipient() tele.Recipient                    { return m.chat }
func (m *addedToGroupCtx) Text() string                                 { return "" }
func (m *addedToGroupCtx) ThreadID() int                                { return 0 }
func (m *addedToGroupCtx) Entities() tele.Entities                      { return nil }
func (m *addedToGroupCtx) Data() string                                 { return "" }
func (m *addedToGroupCtx) Args() []string                               { return nil }
func (m *addedToGroupCtx) Send(what interface{}, opts ...interface{}) error {
	m.sent = what.(string)
	return nil
}
func (m *addedToGroupCtx) SendAlbum(a tele.Album, opts ...interface{}) error       { return nil }
func (m *addedToGroupCtx) Reply(what interface{}, opts ...interface{}) error       { return nil }
func (m *addedToGroupCtx) Forward(msg tele.Editable, opts ...interface{}) error    { return nil }
func (m *addedToGroupCtx) ForwardTo(to tele.Recipient, opts ...interface{}) error  { return nil }
func (m *addedToGroupCtx) Edit(what interface{}, opts ...interface{}) error        { return nil }
func (m *addedToGroupCtx) EditCaption(caption string, opts ...interface{}) error   { return nil }
func (m *addedToGroupCtx) EditOrSend(what interface{}, opts ...interface{}) error  { return nil }
func (m *addedToGroupCtx) EditOrReply(what interface{}, opts ...interface{}) error { return nil }
func (m *addedToGroupCtx) Delete() error                                           { return nil }
func (m *addedToGroupCtx) DeleteAfter(d time.Duration) *time.Timer                 { return nil }
func (m *addedToGroupCtx) Notify(action tele.ChatAction) error                     { return nil }
func (m *addedToGroupCtx) Ship(what ...interface{}) error                          { return nil }
func (m *addedToGroupCtx) Accept(errorMessage ...string) error                     { return nil }
func (m *addedToGroupCtx) Answer(resp *tele.QueryResponse) error                   { return nil }
func (m *addedToGroupCtx) AnswerGuest(result tele.Result) error                    { return nil }
func (m *addedToGroupCtx) Respond(resp ...*tele.CallbackResponse) error            { return nil }
func (m *addedToGroupCtx) RespondText(text string) error                           { return nil }
func (m *addedToGroupCtx) RespondAlert(text string) error                          { return nil }
func (m *addedToGroupCtx) Get(key string) interface{}                              { return nil }
func (m *addedToGroupCtx) Set(key string, val interface{})                         {}

func TestHandleAddedToGroupAdmin(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	const adminID int64 = 123
	cfg := Config{Admins: []int64{adminID}}
	bot := newTestBotWithConfig(dbMock, cfg)

	ctx := &addedToGroupCtx{
		sender: &tele.User{ID: adminID},
		chat:   &tele.Chat{ID: -1, Title: "test", Type: "group"},
	}

	err := bot.handleAddedToGroup(ctx)
	assert.NoError(t, err)
	assert.Empty(t, ctx.sent, "admin should not receive rejection")
	channels := bot.channels.Channels()
	assert.Contains(t, channels, Channel{ID: -1, Title: "test"}, "channel should be tracked")
}

func TestHandleAddedToGroupNonAdmin(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	cfg := Config{Admins: []int64{123}}
	bot := newTestBotWithConfig(dbMock, cfg)

	ctx := &addedToGroupCtx{
		sender: &tele.User{ID: 456, Username: "admin"},
		chat:   &tele.Chat{ID: -1, Title: "test", Type: "group"},
	}

	err := bot.handleAddedToGroup(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "Only admins can add me", ctx.sent)
	assert.Empty(t, bot.channels.Channels(), "channel should not be tracked for non-admin")
}

func TestHandleAddedToGroupZeroID(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBotWithConfig(dbMock, Config{Admins: []int64{0}})

	ctx := &addedToGroupCtx{
		sender: &tele.User{Username: "admin"},
		chat:   &tele.Chat{ID: -2, Title: "test2", Type: "supergroup"},
	}

	err := bot.handleAddedToGroup(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "Only admins can add me", ctx.sent)
	assert.Empty(t, bot.channels.Channels(), "channel should not be tracked for a sender without an ID")
}

func TestHandleAddedToGroupAdminWithoutUsername(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	const adminID int64 = 123
	cfg := Config{Admins: []int64{adminID}}
	bot := newTestBotWithConfig(dbMock, cfg)

	ctx := &addedToGroupCtx{
		sender: &tele.User{ID: adminID, FirstName: "adminuser"},
		chat:   &tele.Chat{ID: -2, Title: "test2", Type: "supergroup"},
	}

	err := bot.handleAddedToGroup(ctx)
	assert.NoError(t, err)
	assert.Empty(t, ctx.sent, "admin should not receive rejection")
	assert.Contains(t, bot.channels.Channels(), Channel{ID: -2, Title: "test2"}, "channel should be tracked")
}
