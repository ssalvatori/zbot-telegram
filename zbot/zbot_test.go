package zbot

import (
	"errors"
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

// stubAPI embeds tele.API so only the methods we need are implemented.
type stubAPI struct {
	tele.API
	leftChat *tele.Chat
	leaveErr error
	sentTo   []int64
	messages []string
	sendErr  error
}

func (s *stubAPI) Leave(chat tele.Recipient) error {
	if s.leaveErr != nil {
		return s.leaveErr
	}
	if c, ok := chat.(*tele.Chat); ok {
		s.leftChat = c
	}
	return nil
}

func (s *stubAPI) Send(to tele.Recipient, what interface{}, opts ...interface{}) (*tele.Message, error) {
	if s.sendErr != nil {
		return nil, s.sendErr
	}
	if u, ok := to.(*tele.User); ok {
		s.sentTo = append(s.sentTo, u.ID)
	}
	if str, ok := what.(string); ok {
		s.messages = append(s.messages, str)
	}
	return &tele.Message{ID: len(s.messages)}, nil
}

// addedToGroupCtx implements tele.Context for testing handleAddedToGroup.
type addedToGroupCtx struct {
	sender *tele.User
	chat   *tele.Chat
	msg    *tele.Message
	sent   string
	api    tele.API
}

func (m *addedToGroupCtx) Bot() tele.API                                { return m.api }
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

	api := &stubAPI{}
	ctx := &addedToGroupCtx{
		sender: &tele.User{ID: adminID},
		chat:   &tele.Chat{ID: -1, Title: "test", Type: "group"},
		api:    api,
	}

	err := bot.handleAddedToGroup(ctx)
	assert.NoError(t, err)
	assert.Empty(t, ctx.sent, "admin should not receive rejection")
	assert.Nil(t, api.leftChat, "bot should not leave when added by an admin")
	channels := bot.channels.Channels()
	assert.Contains(t, channels, Channel{ID: -1, Title: "test"}, "channel should be tracked")
}

func TestHandleAddedToGroupNonAdmin(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	cfg := Config{Admins: []int64{123}}
	bot := newTestBotWithConfig(dbMock, cfg)

	api := &stubAPI{}
	ctx := &addedToGroupCtx{
		sender: &tele.User{ID: 456, Username: "admin"},
		chat:   &tele.Chat{ID: -1, Title: "test", Type: "group"},
		api:    api,
	}

	err := bot.handleAddedToGroup(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "Only admins can add me", ctx.sent)
	assert.Equal(t, int64(-1), api.leftChat.ID, "bot should leave the chat for non-admin inviter")
	assert.Empty(t, bot.channels.Channels(), "channel should not be tracked for non-admin")
}

func TestHandleAddedToGroupZeroID(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBotWithConfig(dbMock, Config{Admins: []int64{0}})

	api := &stubAPI{}
	ctx := &addedToGroupCtx{
		sender: &tele.User{Username: "admin"},
		chat:   &tele.Chat{ID: -2, Title: "test2", Type: "supergroup"},
		api:    api,
	}

	err := bot.handleAddedToGroup(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "Only admins can add me", ctx.sent)
	assert.Equal(t, int64(-2), api.leftChat.ID, "bot should leave the chat when sender has no ID")
	assert.Empty(t, bot.channels.Channels(), "channel should not be tracked for a sender without an ID")
}

func TestHandleAddedToGroupAdminWithoutUsername(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	const adminID int64 = 123
	cfg := Config{Admins: []int64{adminID}}
	bot := newTestBotWithConfig(dbMock, cfg)

	api := &stubAPI{}
	ctx := &addedToGroupCtx{
		sender: &tele.User{ID: adminID, FirstName: "adminuser"},
		chat:   &tele.Chat{ID: -2, Title: "test2", Type: "supergroup"},
		api:    api,
	}

	err := bot.handleAddedToGroup(ctx)
	assert.NoError(t, err)
	assert.Empty(t, ctx.sent, "admin should not receive rejection")
	assert.Nil(t, api.leftChat, "bot should not leave when added by an admin")
	assert.Contains(t, bot.channels.Channels(), Channel{ID: -2, Title: "test2"}, "channel should be tracked")
}

// myChatMemberCtx implements tele.Context for testing handleMyChatMember.
type myChatMemberCtx struct {
	member *tele.ChatMemberUpdate
}

func (m *myChatMemberCtx) Bot() tele.API                                { return nil }
func (m *myChatMemberCtx) Update() tele.Update                          { return tele.Update{} }
func (m *myChatMemberCtx) Message() *tele.Message                       { return nil }
func (m *myChatMemberCtx) Callback() *tele.Callback                     { return nil }
func (m *myChatMemberCtx) Query() *tele.Query                           { return nil }
func (m *myChatMemberCtx) InlineResult() *tele.InlineResult             { return nil }
func (m *myChatMemberCtx) ShippingQuery() *tele.ShippingQuery           { return nil }
func (m *myChatMemberCtx) PreCheckoutQuery() *tele.PreCheckoutQuery     { return nil }
func (m *myChatMemberCtx) Payment() *tele.Payment                       { return nil }
func (m *myChatMemberCtx) Poll() *tele.Poll                             { return nil }
func (m *myChatMemberCtx) PollAnswer() *tele.PollAnswer                 { return nil }
func (m *myChatMemberCtx) ChatMember() *tele.ChatMemberUpdate           { return m.member }
func (m *myChatMemberCtx) ChatJoinRequest() *tele.ChatJoinRequest       { return nil }
func (m *myChatMemberCtx) Migration() (int64, int64)                    { return 0, 0 }
func (m *myChatMemberCtx) Topic() *tele.Topic                           { return nil }
func (m *myChatMemberCtx) Boost() *tele.BoostUpdated                    { return nil }
func (m *myChatMemberCtx) BoostRemoved() *tele.BoostRemoved             { return nil }
func (m *myChatMemberCtx) PurchasedPaidMedia() *tele.PaidMediaPurchased { return nil }
func (m *myChatMemberCtx) Sender() *tele.User {
	if m.member != nil {
		return m.member.Sender
	}
	return nil
}
func (m *myChatMemberCtx) Chat() *tele.Chat {
	if m.member != nil {
		return m.member.Chat
	}
	return nil
}
func (m *myChatMemberCtx) Recipient() tele.Recipient { return m.Chat() }
func (m *myChatMemberCtx) Text() string              { return "" }
func (m *myChatMemberCtx) ThreadID() int             { return 0 }
func (m *myChatMemberCtx) Entities() tele.Entities   { return nil }
func (m *myChatMemberCtx) Data() string              { return "" }
func (m *myChatMemberCtx) Args() []string            { return nil }
func (m *myChatMemberCtx) Send(what interface{}, opts ...interface{}) error {
	return nil
}
func (m *myChatMemberCtx) SendAlbum(a tele.Album, opts ...interface{}) error       { return nil }
func (m *myChatMemberCtx) Reply(what interface{}, opts ...interface{}) error       { return nil }
func (m *myChatMemberCtx) Forward(msg tele.Editable, opts ...interface{}) error    { return nil }
func (m *myChatMemberCtx) ForwardTo(to tele.Recipient, opts ...interface{}) error  { return nil }
func (m *myChatMemberCtx) Edit(what interface{}, opts ...interface{}) error        { return nil }
func (m *myChatMemberCtx) EditCaption(caption string, opts ...interface{}) error   { return nil }
func (m *myChatMemberCtx) EditOrSend(what interface{}, opts ...interface{}) error  { return nil }
func (m *myChatMemberCtx) EditOrReply(what interface{}, opts ...interface{}) error { return nil }
func (m *myChatMemberCtx) Delete() error                                           { return nil }
func (m *myChatMemberCtx) DeleteAfter(d time.Duration) *time.Timer                 { return nil }
func (m *myChatMemberCtx) Notify(action tele.ChatAction) error                     { return nil }
func (m *myChatMemberCtx) Ship(what ...interface{}) error                          { return nil }
func (m *myChatMemberCtx) Accept(errorMessage ...string) error                     { return nil }
func (m *myChatMemberCtx) Answer(resp *tele.QueryResponse) error                   { return nil }
func (m *myChatMemberCtx) AnswerGuest(result tele.Result) error                    { return nil }
func (m *myChatMemberCtx) Respond(resp ...*tele.CallbackResponse) error            { return nil }
func (m *myChatMemberCtx) RespondText(text string) error                           { return nil }
func (m *myChatMemberCtx) RespondAlert(text string) error                          { return nil }
func (m *myChatMemberCtx) Get(key string) interface{}                              { return nil }
func (m *myChatMemberCtx) Set(key string, val interface{})                         {}

func TestHandleMyChatMemberKicked(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBotWithConfig(dbMock, Config{Channels: []Channel{{ID: -100, Title: "gone"}}})

	ctx := &myChatMemberCtx{
		member: &tele.ChatMemberUpdate{
			Chat:          &tele.Chat{ID: -100, Title: "gone", Type: "supergroup"},
			Sender:        &tele.User{ID: 1, Username: "admin"},
			NewChatMember: &tele.ChatMember{Role: tele.Kicked},
		},
	}

	err := bot.handleMyChatMember(ctx)
	assert.NoError(t, err)
	assert.NotContains(t, bot.channels.Channels(), Channel{ID: -100, Title: "gone"}, "channel should be forgotten")
}

func TestHandleMyChatMemberLeft(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBotWithConfig(dbMock, Config{Channels: []Channel{{ID: -200, Title: "left"}}})

	ctx := &myChatMemberCtx{
		member: &tele.ChatMemberUpdate{
			Chat:          &tele.Chat{ID: -200, Title: "left", Type: "group"},
			Sender:        &tele.User{ID: 2},
			NewChatMember: &tele.ChatMember{Role: tele.Left},
		},
	}

	err := bot.handleMyChatMember(ctx)
	assert.NoError(t, err)
	assert.Empty(t, bot.channels.Channels(), "channel should be forgotten on left")
}

func TestHandleMyChatMemberStillMember(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBotWithConfig(dbMock, Config{Channels: []Channel{{ID: -300, Title: "stays"}}})

	ctx := &myChatMemberCtx{
		member: &tele.ChatMemberUpdate{
			Chat:          &tele.Chat{ID: -300, Title: "stays", Type: "supergroup"},
			Sender:        &tele.User{ID: 3},
			NewChatMember: &tele.ChatMember{Role: tele.Member},
		},
	}

	err := bot.handleMyChatMember(ctx)
	assert.NoError(t, err)
	assert.Contains(t, bot.channels.Channels(), Channel{ID: -300, Title: "stays"}, "channel should remain when bot is still a member")
}

func TestHandleMyChatMemberNil(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBot(dbMock)

	ctx := &myChatMemberCtx{member: nil}

	err := bot.handleMyChatMember(ctx)
	assert.NoError(t, err, "nil member update should be a no-op")
}

func TestNotifyAdmins(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	cfg := Config{Admins: []int64{111, 222}}
	bot := newTestBotWithConfig(dbMock, cfg)

	api := &stubAPI{}
	ctx := &addedToGroupCtx{api: api}

	bot.NotifyAdmins(ctx, "hello admins")

	assert.Equal(t, []int64{111, 222}, api.sentTo, "message should be sent to every admin")
	assert.Equal(t, []string{"hello admins", "hello admins"}, api.messages)
}

func TestNotifyAdminsSendError(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBotWithConfig(dbMock, Config{Admins: []int64{111}})

	api := &stubAPI{sendErr: errors.New("boom")}
	ctx := &addedToGroupCtx{api: api}

	bot.NotifyAdmins(ctx, "msg") // send error is logged, not propagated
	assert.Empty(t, api.sentTo)
}

func TestEnableOrAddChannelNew(t *testing.T) {
	exists := false
	dbMock := &db.ZbotDatabaseMock{ChannelExists: &exists}
	bot := newTestBot(dbMock)

	err := bot.EnableOrAddChannel(tele.Chat{ID: -10, Title: "newgroup"}, tele.User{ID: 1, Username: "admin"})
	assert.NoError(t, err)
	assert.True(t, dbMock.AddedChannel, "AddChannel should be called for a new channel")
	assert.False(t, dbMock.EnabledChannel, "EnableChannel should not be called for a new channel")
	assert.Contains(t, bot.channels.Channels(), Channel{ID: -10, Title: "newgroup"}, "channel should be tracked")
}

func TestEnableOrAddChannelExisting(t *testing.T) {
	exists := true
	dbMock := &db.ZbotDatabaseMock{ChannelExists: &exists}
	bot := newTestBot(dbMock)

	err := bot.EnableOrAddChannel(tele.Chat{ID: -20, Title: "oldgroup"}, tele.User{ID: 2, Username: "admin"})
	assert.NoError(t, err)
	assert.True(t, dbMock.EnabledChannel, "EnableChannel should be called for an existing channel")
	assert.False(t, dbMock.AddedChannel, "AddChannel should not be called for an existing channel")
	assert.Contains(t, bot.channels.Channels(), Channel{ID: -20, Title: "oldgroup"}, "channel should be tracked")
}

func TestRejectAndLeave(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBot(dbMock)

	api := &stubAPI{}
	ctx := &addedToGroupCtx{chat: &tele.Chat{ID: -5, Title: "bad"}, api: api}

	err := bot.rejectAndLeave(ctx, "nope")
	assert.NoError(t, err)
	assert.Equal(t, "nope", ctx.sent, "rejection message should be sent")
	assert.Equal(t, int64(-5), api.leftChat.ID, "bot should leave the chat")
}

func TestRejectAndLeaveError(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBot(dbMock)

	api := &stubAPI{leaveErr: errors.New("cannot leave")}
	ctx := &addedToGroupCtx{chat: &tele.Chat{ID: -6}, api: api}

	err := bot.rejectAndLeave(ctx, "nope")
	assert.Error(t, err, "leave error should be propagated")
	assert.Equal(t, "nope", ctx.sent, "rejection message is still sent before leaving")
}

func TestRejectAndLeaveNoBot(t *testing.T) {
	dbMock := &db.ZbotDatabaseMock{}
	bot := newTestBot(dbMock)

	ctx := &addedToGroupCtx{chat: &tele.Chat{ID: -7}, api: nil}

	err := bot.rejectAndLeave(ctx, "nope")
	assert.NoError(t, err, "missing bot in context is not an error")
	assert.Equal(t, "nope", ctx.sent)
}
