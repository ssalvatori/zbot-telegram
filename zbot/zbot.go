package zbot

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	command "github.com/ssalvatori/zbot-telegram/commands"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
	zbotMiddleware "github.com/ssalvatori/zbot-telegram/zbot/middleware"

	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

// ExternalModule definition
type ExternalModule struct {
	Key         string
	File        string
	Description string
}

// Flags feature toggles; false means the feature is disabled
type Flags struct {
	Ignore bool
	Level  bool
}

// Config is everything the bot needs to run. Build it with ConfigFromFile or
// fill it directly in tests.
type Config struct {
	// Token Telegram API Token (key:secret Format)
	Token string
	// ModulesPath Absolute path where the modules are located
	ModulesPath string
	// IgnoreDuration Ignore a user for this amount of seconds
	IgnoreDuration int
	Flags          Flags
	// Modules list of extra modules
	Modules []ExternalModule
	// Channels pre-configured channels (the list grows as messages arrive)
	Channels []Channel
	// Levels minimum user level per command
	Levels command.Levels
	// List of global admins.
	Admins []int64
	// LogFileName is the name of the file where logs are written in addition to stdout.
	LogFileName string
}

// Version information, injected at build time with -ldflags -X
var (
	version   = "dev-master"
	buildTime = time.Now().Format("2006-01-02 15:04:05")
	gitHash   = "undefined"
)

// DefaultLevels minimum user level required per command when Flags.Level is on
var DefaultLevels = command.Levels{
	Ignore:   100,
	Lock:     1000,
	Learn:    0,
	Append:   0,
	Forget:   1000,
	Who:      0,
	Top:      0,
	Stats:    0,
	Version:  0,
	Ping:     0,
	Last:     0,
	Rand:     0,
	Find:     0,
	Get:      0,
	Search:   0,
	External: 0,
	Level:    0,
	Help:     0,
}

const DEFAULT_IGNORE_DURATION = 300

// commandPrefix matches messages that should be processed as bot commands (starting with ! or ?)
var commandPrefix = regexp.MustCompile(`^[!?]`)

// LoadConfig maps the on-disk configuration into a Config.
func LoadConfig(conf *Configuration) (Config, error) {
	cfg := Config{
		Token:          conf.Zbot.Token,
		ModulesPath:    conf.Modules.Path,
		IgnoreDuration: conf.Zbot.IgnoreDuration, //TODO: validate if this value is > 0, otherwise set to default DEFAULT_IGNORE_DURATION
		Flags:          Flags{Ignore: conf.Zbot.Ignore, Level: conf.Zbot.Level},
		Levels:         DefaultLevels,
		Admins:         conf.Zbot.Admins,
	}

	if err := mapstructure.Decode(conf.Modules.List, &cfg.Modules); err != nil {
		return Config{}, fmt.Errorf("decode modules: %w", err)
	}

	return cfg, nil
}

// Bot is a running zbot instance.
type Bot struct {
	cfg       Config                   // Configuration, loaded from file or built in tests
	db        db.ZbotDatabase          // Database connection, initialized in Run
	commands  *command.CommandsList    // List of commands the bot can run, including built-in and learned commands (chained list)
	adminCmds []command.CommandElement // List of commands the bot can run for admins
	channels  *channelRegistry         // channels the bot is in, updated as invitations arrive from Bot Administrator
	admins    []Admin                  // List of bot administrators, loaded from configuration and verified at startup
}

// New wires a Bot together. It does not touch the network or the database;
// that happens in Run.
func New(cfg Config, database db.ZbotDatabase) (*Bot, error) {
	if database == nil {
		return nil, errors.New("database not configured")
	}
	if cfg.ModulesPath == "" {
		cfg.ModulesPath = utils.GetCurrentDirectory() + "/../modules/"
	}

	//TODO:: this can be remove when the validation is done in while loading the configuration file.
	if cfg.IgnoreDuration <= 0 {
		cfg.IgnoreDuration = DEFAULT_IGNORE_DURATION
	}

	// The command package still keeps these as globals; this is the only place we set them.
	command.Setup()

	b := &Bot{
		cfg:      cfg,
		db:       database,
		channels: newChannelRegistry(cfg.Channels),
	}

	b.commands = b.buildBaseCommands()
	b.commands.List.PushBackList(b.buildLearnCommands().List)

	b.adminCmds = b.buildAdminCommands()

	return b, nil
}

// buildLearnCommand creates the chain of learn commands once.
func (b *Bot) buildLearnCommands() *command.CommandsList {

	lv := b.cfg.Levels
	d := b.db

	cl := &command.CommandsList{List: list.New()}
	cl.Chain("last", &command.LastCommand{Db: b.db}, lv.Last)
	cl.Chain("rand", &command.RandCommand{Db: b.db}, lv.Rand)
	cl.Chain("who", &command.WhoCommand{Db: b.db}, lv.Who)
	cl.Chain("find", &command.FindCommand{Db: d}, lv.Find)
	cl.Chain("get", &command.GetCommand{Db: d}, lv.Get)
	cl.Chain("search", &command.SearchCommand{Db: d}, lv.Search)
	cl.Chain("stats", &command.StatsCommand{Db: d}, lv.Stats)
	cl.Chain("learn", &command.LearnCommand{Db: d}, lv.Learn)
	cl.Chain("append", &command.AppendCommand{Db: d}, lv.Append)
	cl.Chain("forget", &command.ForgetCommand{Db: d}, lv.Forget)
	cl.Chain("lock", &command.LockCommand{Db: d}, lv.Lock)

	return cl
}

// buildCommands creates the chain of built-in commands once.
func (b *Bot) buildBaseCommands() *command.CommandsList {
	lv := b.cfg.Levels
	d := b.db

	cl := &command.CommandsList{List: list.New()}
	cl.Chain("help", &command.HelpCommand{Commands: cl}, lv.Help)
	cl.Chain("ping", &command.PingCommand{}, lv.Ping)
	cl.Chain("version", &command.VersionCommand{GitHash: gitHash, Version: version, BuildTime: buildTime}, lv.Version)
	cl.Chain("level", &command.LevelCommand{Db: d}, lv.Level)
	cl.Chain("ignore", &command.IgnoreCommand{Db: d, IgnoreDuration: b.cfg.IgnoreDuration, AdminLevel: lv.Ignore}, lv.Ignore)
	return cl
}

// buildAdminCommands list of commands that can be run by admins in private chat with the bot.
func (b *Bot) buildAdminCommands() []command.CommandElement {
	cmds := []command.CommandElement{
		*command.NewCommandElement("/version", 0, &command.VersionCommand{GitHash: gitHash, Version: version, BuildTime: buildTime}),
		*command.NewCommandElement("/ping", 0, &command.PingCommand{}),
		*command.NewCommandElement("/logs", 0, &command.LogsCommand{Path: command.DefaultLogsPath}),
		// command.NewCommandElement("/help", 0, "/help", command.LevelHelpMessage{}),
	}
	return cmds
}

// Run connects to Telegram and the database and processes updates until ctx
// is cancelled or a fatal error occurs.
func (b *Bot) Run(ctx context.Context) error {
	slog.Info("Loading zbot-telegram", "version", version, "buildTime", buildTime, "gitHash", gitHash)
	slog.Info("Database", "engine", b.db.Engine(), "connection", b.db.GetConnectionInfo())
	slog.Info("Modules", "path", b.cfg.ModulesPath)
	slog.Info("Configuration Flags", "ignore", b.cfg.Flags.Ignore, "level", b.cfg.Flags.Level)

	poller := &tele.LongPoller{Timeout: 10 * time.Second}

	tb, err := tele.NewBot(tele.Settings{
		Token:  b.cfg.Token,
		Poller: tele.NewMiddlewarePoller(poller, b.filterUpdate),
	})
	if err != nil {
		return fmt.Errorf("bot init: %w", err)
	}
	tb.Use(middleware.Logger())
	// tb.Use(middleware.RegistryChannel()) #TODO: Create a middleware to register the channel the bot is in.

	if err := b.db.Init(); err != nil {
		return fmt.Errorf("db init: %w", err)
	}
	defer b.db.Close()

	if b.cfg.Flags.Ignore {
		if err := b.db.UserCleanupIgnorelist(); err != nil {
			slog.Error("ignore list cleanup error", "err", err)
		}
	}

	b.admins = b.setupAdministrators(tb)

	b.cleanUpCommands(tb)
	b.registerAdminCommands(tb)

	// b.registerExternalModules(tb)
	// tb.Handle(tele.OnText, b.handleText)
	tb.Handle(tele.OnAddedToGroup, b.handleAddedToGroup)
	tb.Handle(tele.OnMyChatMember, b.handleMyChatMember)

	// adminGroup := tb.Group()
	// adminGroup.Use(zbotMiddleware.AdminOnly(b.IsAdmin))

	done := make(chan struct{})
	go func() {
		defer close(done)
		tb.Start()
	}()

	<-ctx.Done()
	slog.Info("Shutting down")
	tb.Stop()
	<-done
	return nil
}

// cleanUpCommands removes all commands from the bot, so that we can set them again with the current configuration.
func (b *Bot) cleanUpCommands(tb *tele.Bot) {
	slog.Debug("Cleaning up commands for all admins", "admin_ids", b.cfg.Admins)

	if err := tb.DeleteCommands(); err != nil {
		slog.Error("delete commands error", "err", err)
	}

}

// registerExternalModules installs a /<key> handler for every configured module whose binary exists.
func (b *Bot) registerExternalModules(tb *tele.Bot) {
	slog.Debug("Modules to load", "modules", b.cfg.Modules)
	botCommands := []tele.Command{}

	for _, module := range b.cfg.Modules {
		cmdString := "/" + module.Key
		slog.Debug("Loading module", "key", module.Key, "path", b.cfg.ModulesPath+module.File)

		if _, err := command.LookPathCommand(b.cfg.ModulesPath + module.File); err != nil {
			slog.Error("module file not found", "file", module.File, "key", module.Key, "path", b.cfg.ModulesPath)
			continue
		}

		tb.Handle(cmdString, b.handleExternalModule)
		botCommands = append(botCommands, tele.Command{Text: cmdString, Description: module.Description})
	}

	slog.Debug("Setting bot commands", "commands", botCommands)
	if err := tb.SetCommands(botCommands); err != nil {
		slog.Error("set commands error", "err", err)
	}
}

func (b *Bot) handleExternalModule(c tele.Context) error {
	slog.Debug("incoming message", "user", c.Sender().Username, "text", c.Text(), "payload", c.Message().Payload, "private", c.Message().Private())

	response := b.runExternalModule(c.Message())
	if err := c.Send(response); err != nil {
		slog.Error("send error", "err", err)
		return err
	}
	return nil
}

// handleText is the main entry point for processing messages.
// It decides whether to ignore the message, run a command, or do nothing.
func (b *Bot) handleText(c tele.Context) error {

	slog.Debug("processing message from", "user", c.Sender())

	if c.Message().Private() {
		return zbotMiddleware.AdminOnly(b.IsAdmin)(b.handleAdminText)(c)
	}

	response := b.messagesProcessing(c.Message())
	if response == "" {
		return nil
	}

	if err := c.Send(response); err != nil {
		slog.Error("send response error", "err", err)
		return err
	}
	return nil
}

// handleAdminText processes private-chat messages; only reachable by admins via AdminOnly.
func (b *Bot) handleAdminText(c tele.Context) error {
	slog.Debug("admin message received", "user", c.Sender().Username)

	// TODO: process messages using b.adminCommands
	return nil
}

// handleAddedToGroup processes the bot being added to a group.
// Only admins may add the bot; non-admin inviters get a rejection message
// and the bot leaves the chat immediately.
func (b *Bot) handleAddedToGroup(c tele.Context) error {
	sender := c.Sender()
	if sender == nil || sender.ID == 0 {
		slog.Warn("inviter has no Telegram ID, cannot verify admin status")
		return b.rejectAndLeave(c, "Only admins can add me")
	}

	if !b.IsAdmin(sender.ID) {
		slog.Info("non-admin tried to add the bot", "user_id", sender.ID, "chat_id", c.Chat().ID)
		return b.rejectAndLeave(c, "Only admins can add me")
	}

	b.EnableOrAddChannel(*c.Chat(), *c.Sender()) // Add the channel to the database and the in-memory registry
	b.NotifyAdmins(c, fmt.Sprintf("Bot added to new group: %s (%d) by @%s (%d)", c.Chat().Title, c.Chat().ID, sender.Username, sender.ID))

	return nil
}

// NotifyAdmins sends a message to all configured zbot admins.
func (b *Bot) NotifyAdmins(c tele.Context, message string) {
	slog.Debug("Notify admins", "message", message, "admin_ids", b.cfg.Admins)

	for _, adminID := range b.cfg.Admins {
		admin := &tele.User{ID: adminID}
		msg, err := c.Bot().Send(admin, message)
		if err != nil {
			slog.Error("notify admin error", "admin_id", adminID, "err", err)
			continue
		}
		slog.Debug("notified admin", "admin_id", adminID, "message_id", msg.ID)
	}
}

// EnableOrAddChannel adds a channel to the database and the in-memory registry.
func (b *Bot) EnableOrAddChannel(chat tele.Chat, sender tele.User) error {
	// Add the channel to the DB or activate it if it already exists
	// Send messages to the admins that the bot has been added to a new group and inform the inviter.

	slog.Info("Checking if channel should be enabled or added", "chat", chat.Title, "chat_id", chat.ID, "inviter", sender.Username, "inviter_id", sender.ID)

	if !b.db.CheckIfChannelExist(chat.ID) {
		slog.Info("adding new channel to database", "chat", chat.Title, "id", chat.ID)
		b.db.AddChannel(chat.ID, chat.Title)
	} else {
		slog.Info("enabling existing channel in database", "chat", chat.Title, "id", chat.ID)
		b.db.EnableChannel(chat.ID)
	}

	b.channels.Track(chat)

	return nil
}

// rejectAndLeave sends a rejection message and leaves the chat.
func (b *Bot) rejectAndLeave(c tele.Context, message string) error {
	if err := c.Send(message); err != nil {
		slog.Error("send rejection error", "err", err)
	}

	bot := c.Bot()
	if bot == nil {
		slog.Warn("cannot leave chat: no bot in context")
		return nil
	}

	if err := bot.Leave(c.Chat()); err != nil {
		slog.Error("leave chat error", "chat_id", c.Chat().ID, "err", err)
		return err
	}

	slog.Info("left unauthorized chat", "chat", c.Chat().Title, "chat_id", c.Chat().ID)
	return nil
}

// handleMyChatMember reacts to changes in the bot's own membership status.
// When the bot is kicked or leaves a chat, the channel is removed from the
// database and from the in-memory registry.
func (b *Bot) handleMyChatMember(c tele.Context) error {
	m := c.ChatMember()
	if m == nil || m.NewChatMember == nil || m.Chat == nil {
		return nil
	}

	role := m.NewChatMember.Role
	if role != tele.Kicked && role != tele.Left {
		return nil
	}

	slog.Info("bot removed from chat", "chat", m.Chat.Title, "chat_id", m.Chat.ID, "status", role)

	removed, err := b.db.DisableChannel(m.Chat.ID)
	if err != nil {
		slog.Error("remove channel from db error", "chat_id", m.Chat.ID, "err", err)
	} else if removed {
		slog.Debug("channel removed from db", "chat_id", m.Chat.ID)
	}

	b.channels.Forget(m.Chat.ID)
	return nil
}

// filterUpdate is the poller middleware: drops spam and records the chats the bot is in.
func (b *Bot) filterUpdate(u *tele.Update) bool {
	if u.Message == nil {
		return true
	}

	if strings.Contains(u.Message.Text, "spam") {
		return false
	}

	if u.Message.Chat != nil && (u.Message.Chat.Type == "group" || u.Message.Chat.Type == "supergroup") {
		b.channels.Track(*u.Message.Chat)
	}

	return true
}

func (b *Bot) runExternalModule(message *tele.Message) string {

	cmd, err := utils.ParseCommand(message.Text)
	if err != nil {
		slog.Error("parse command error", "err", err)
		return ""
	}

	cmdFile, err := utils.GetCommandFile(cmd, b.cfg.Modules)
	if err != nil {
		slog.Error("get command file error", "err", err)
		return ""
	}

	fullPathToBinary, _ := command.LookPathCommand(b.cfg.ModulesPath + cmdFile)

	chatName := ""
	if message.Chat.Type != "private" {
		chatName = message.Chat.Title
	}

	u := user.BuildUser(message.Sender, b.db)
	slog.Debug("Running module", "path", fullPathToBinary)
	return utils.RunExternalCommand(fullPathToBinary, u.Username, strconv.Itoa(u.Level), chatName, strings.TrimSpace(message.Payload))
}

// messagesProcessing decides whether a message is a command for the bot and runs it.
func (b *Bot) messagesProcessing(message *tele.Message) string {

	slog.Debug("message", "chat", message, "text", message.Text)

	private := message.Chat.Type == "private"
	username := strings.ToLower(message.Sender.Username)

	channel, err := b.resolveChannel(message.Chat)
	if err != nil {
		slog.Error("resolve channel error", "chat", message.Chat.ID, "title", message.Chat.Title, "err", err)
		return ""
	}

	if b.ignoreMessage(username, channel.ID) {
		slog.Debug("message ignored", "user", username)
		return ""
	}

	if !commandPrefix.MatchString(message.Text) {
		return ""
	}

	slog.Debug("received command", "user", username, "text", message.Text)
	return b.cmdProcessing(*message, channel, private)
}

// resolveChannel looks up the stored channel for a Telegram chat by its ID.
func (b *Bot) resolveChannel(chat *tele.Chat) (db.Channel, error) {
	if chat == nil {
		return db.Channel{}, errors.New("message has no chat")
	}
	return b.db.GetChannel(chat.ID)
}

// ignoreMessage check if the message should be ignore base on the username and channel, if the ignore flag is disable we don't ignore any messages.
func (b *Bot) ignoreMessage(username string, channelID uint) bool {
	if !b.cfg.Flags.Ignore {
		slog.Debug("ignore flag is disabled, not ignoring messages")
		return false
	}
	return b.db.UserCheckIgnore(username, channelID)
}

// cmdProcessing runs a command message through the command chain.
func (b *Bot) cmdProcessing(msg tele.Message, channel db.Channel, private bool) string {

	commandName := command.GetCommandInformation(msg.Text)

	if command.IsCommandDisabled(commandName) {
		slog.Debug("command disabled", "command", commandName)
		return fmt.Sprintf("Command [%v] is disabled", commandName)
	}

	u := user.BuildUser(msg.Sender, b.db)

	if b.cfg.Flags.Level {
		requiredLevel := command.GetMinimumLevel(commandName, b.cfg.Levels)
		if !command.CheckPermission(commandName, u, requiredLevel) {
			return fmt.Sprintf("Your level is not enough < %d", requiredLevel)
		}
	}

	messageString := msg.Text
	if msg.ReplyTo != nil {
		messageString = fmt.Sprintf("%s %s %s", messageString, msg.ReplyTo.Sender.Username, msg.ReplyTo.Text)
	}

	return b.commands.Run(commandName, messageString, u, int(channel.ID), channel.Name, private)
}
