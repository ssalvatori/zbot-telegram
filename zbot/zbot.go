package zbot

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	command "github.com/ssalvatori/zbot-telegram/commands"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"

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
	Admins []string
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
	cfg      Config                // Configuration, loaded from file or built in tests
	db       db.ZbotDatabase       // Database connection, initialized in Run
	commands *command.CommandsList // List of commands the bot can run, including built-in and learned commands
	channels *channelRegistry      // channels the bot is in, updated as invitations arrive from Bot Administrator
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

	return b, nil
}

// buildLearnCommand creates the chain of learn commands once.
func (b *Bot) buildLearnCommands() *command.CommandsList {

	lv := b.cfg.Levels
	d := b.db

	cl := &command.CommandsList{List: list.New()}
	cl.Chain("last", &command.LastCommand{Db: b.db}, lv.Last, command.LastHelpMessage)
	cl.Chain("rand", &command.RandCommand{Db: b.db}, lv.Rand, command.RandHelpMessage)
	cl.Chain("who", &command.WhoCommand{Db: b.db}, lv.Who, command.WhoHelpMessage)
	cl.Chain("find", &command.FindCommand{Db: d}, lv.Find, command.FindHelpMessage)
	cl.Chain("get", &command.GetCommand{Db: d}, lv.Get, command.GetHelpMessage)
	cl.Chain("search", &command.SearchCommand{Db: d}, lv.Search, command.SearchHelpMessage)
	cl.Chain("stats", &command.StatsCommand{Db: d}, lv.Stats, command.StatsHelpMessage)
	cl.Chain("learn", &command.LearnCommand{Db: d}, lv.Learn, command.LearnHelpMessage)
	cl.Chain("append", &command.AppendCommand{Db: d}, lv.Append, command.AppendHelpMessage)
	cl.Chain("forget", &command.ForgetCommand{Db: d}, lv.Forget, command.ForgetHelpMessage)
	cl.Chain("lock", &command.LockCommand{Db: d}, lv.Lock, command.LockHelpMessage)

	return cl
}

// buildCommands creates the chain of built-in commands once.
func (b *Bot) buildBaseCommands() *command.CommandsList {
	lv := b.cfg.Levels
	d := b.db

	cl := &command.CommandsList{List: list.New()}
	cl.Chain("help", &command.HelpCommand{Commands: cl}, lv.Help, command.HelpMessage)
	cl.Chain("ping", &command.PingCommand{Db: d}, lv.Ping, command.PingHelpMessage)
	cl.Chain("version", &command.VersionCommand{GitHash: gitHash, Version: version, BuildTime: buildTime}, lv.Version, command.VersionHelpMessage)
	cl.Chain("level", &command.LevelCommand{Db: d}, lv.Level, command.LevelHelpMessage)
	cl.Chain("ignore", &command.IgnoreCommand{Db: d, IgnoreDuration: b.cfg.IgnoreDuration, AdminLevel: lv.Ignore}, lv.Ignore, command.IgnoreHelpMessage)
	return cl
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

	if err := b.db.Init(); err != nil {
		return fmt.Errorf("db init: %w", err)
	}
	defer b.db.Close()

	if b.cfg.Flags.Ignore {
		if err := b.db.UserCleanupIgnorelist(); err != nil {
			slog.Error("ignore list cleanup error", "err", err)
		}
	}

	b.registerExternalModules(tb)
	tb.Handle(tele.OnText, b.handleText)
	tb.Handle(tele.OnAddedToGroup, b.handleAddedToGroup)

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

func (b *Bot) handleText(c tele.Context) error {
	if c.Message().Private() {
		return nil
	}

	response := b.messagesProcessing(c.Message(), c.Chat().Title)
	if response == "" {
		return nil
	}

	if err := c.Send(response); err != nil {
		slog.Error("send response error", "err", err)
		return err
	}
	return nil
}

func (b *Bot) handleAddedToGroup(c tele.Context) error {
	sender := c.Sender()
	inviter := strings.ToLower(sender.Username)

	if inviter == "" {
		slog.Warn("inviter username is empty, cannot verify admin status")
		return c.Send("I cannot verify your admin status because your username is empty. Please set a username in Telegram and try again.")
	}

	if !slices.ContainsFunc(b.cfg.Admins, func(admin string) bool {
		return strings.EqualFold(admin, inviter)
	}) {
		if err := c.Send("Only admins can add me"); err != nil {
			slog.Error("send rejection error", "err", err)
			return err
		}
		return nil
	}

	slog.Debug("added to group", "chat", c.Chat().Title, "id", c.Chat().ID, "inviter", inviter)
	b.channels.Track(*c.Chat())
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
func (b *Bot) messagesProcessing(message *tele.Message, channelName string) string {

	slog.Debug("message", "chat", message.Chat, "text", message.Text)

	private := message.Chat.Type == "private" && channelName == ""
	username := strings.ToLower(message.Sender.Username)

	if b.ignoreMessage(username, channelName) {
		slog.Debug("message ignored", "user", username)
		return ""
	}

	if !commandPrefix.MatchString(message.Text) {
		return ""
	}

	slog.Debug("received command", "user", username, "text", message.Text)
	return b.cmdProcessing(*message, channelName, private)
}

// ignoreMessage check if the message should be ignore base on the username and channelName, if the ignore flag is disable we don't ignore any messages.
func (b *Bot) ignoreMessage(username string, channelName string) bool {
	if !b.cfg.Flags.Ignore {
		slog.Debug("ignore flag is disabled, not ignoring messages")
		return false
	}
	channel, err := b.db.GetChannel(channelName)
	if err != nil {
		slog.Error("get channel error", "err", err)
		return false
	}
	return b.db.UserCheckIgnore(username, channel.ID)
}

// cmdProcessing runs a command message through the command chain.
func (b *Bot) cmdProcessing(msg tele.Message, channelName string, private bool) string {

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

	channel, err := b.db.GetChannel(channelName)

	if err != nil {
		slog.Error("get channel error", "err", err)
		return ""
	}

	return b.commands.Run(commandName, messageString, u, int(channel.ID), channel.Name, private)
}
