package zbot

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"container/list"

	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	command "github.com/ssalvatori/zbot-telegram/commands"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
)

// ExternalModule definition
type ExternalModule struct {
	Key         string
	File        string
	Description string
}

// Channel definition
type Channel struct {
	ID        int64
	Title     string
	AuthToken string
}

// ConfigurationFlags configurations false means the feature is disabled
type ConfigurationFlags struct {
	Ignore bool
	Level  bool
}

// ConfigurationWebhook configuration
type ConfigurationWebhook struct {
	Enable bool
	Port   int
}

var (
	version   = "dev-master"
	buildTime = time.Now().Format("2006-01-02 15:04:05")
	gitHash   = "undefined"
	//DatabaseType database backend to be use (mysql or sqlite)
	DatabaseType = ""
	//APIToken Telegram API Token (key:secret Format)
	APIToken = ""
	//ModulesPath Absolute path where the modules are located
	ModulesPath = utils.GetCurrentDirectory() + "/../modules/"
	//Flags zbot configurations
	Flags = ConfigurationFlags{Ignore: false, Level: false}
	//IgnoreDuration Ignore a user for this amount of seconds
	IgnoreDuration = 300
	//DisableLearnChannels List of channels were Learn modules should be disabled (use comma as separator)
	DisableLearnChannels = ""

	//Channels List of Channels where the bot is present (this list is growing with new messages)
	Channels []Channel

	//ExternalModules List of extra modules
	ExternalModules []ExternalModule

	//Db interface to the database
	Db db.ZbotDatabase

	levelsConfig = command.Levels{
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
	}
)

// Execute run Zbot
func Execute() {
	slog.Info("Loading zbot-telegram", "version", version, "buildTime", buildTime, "gitHash", gitHash)

	slog.Info("Database", "type", DatabaseType)
	slog.Info("Modules", "path", ModulesPath)
	slog.Info("Configuration Flags Ignore", "ignore", Flags.Ignore)
	slog.Info("Configuration Flags Level", "level", Flags.Level)

	command.Setup()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMiddlewares(spamFilterMiddleware),
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(APIToken, opts...)
	if err != nil {
		slog.Error("bot init error", "err", err)
		os.Exit(1)
	}

	err = Db.Init()
	defer Db.Close()

	if err != nil {
		slog.Error("db init error", "err", err)
		os.Exit(1)
	}

	Db.SetIgnoreTime(int64(IgnoreDuration))

	slog.Debug("Modules to load", "modules", ExternalModules)
	botCommands := []models.BotCommand{}

	//Register extra modules
	for _, module := range ExternalModules {
		var cmdString = "/" + module.Key
		slog.Debug("Loading module", "key", module.Key, "path", ModulesPath+module.File)

		_, err := command.LookPathCommand(ModulesPath + module.File)

		if err != nil {
			slog.Error("module file not found", "file", module.File, "key", module.Key, "path", ModulesPath)
			continue
		}

		b.RegisterHandler(bot.HandlerTypeMessageText, cmdString, bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
			msg := update.Message
			if msg == nil {
				return
			}

			slog.Debug("incoming message", "user", msg.From.Username, "text", msg.Text, "private", msg.Chat.Type == "private")

			response := runExternalModule(Db, msg, ExternalModules)
			if response != "" {
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: msg.Chat.ID,
					Text:   response,
				})
				if err != nil {
					slog.Error("send error", "err", err)
				}
			}
		})

		botCommands = append(botCommands, models.BotCommand{Command: module.Key, Description: module.Description})
	}

	slog.Debug("Setting bot commands", "commands", botCommands)
	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: botCommands,
	})
	if err != nil {
		slog.Error("set commands error", "err", err)
	}

	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := update.Message
	if msg == nil {
		return
	}

	// Track channels
	if msg.Chat.Type == "group" || msg.Chat.Type == "supergroup" {
		Channels = appendChannel(Channels, msg.Chat)
	}

	if msg.Chat.Type == "private" {
		return
	}

	chatName := msg.Chat.Title

	var response = messagesProcessing(Db, msg, chatName)
	if response != "" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: msg.Chat.ID,
			Text:   response,
		})
		if err != nil {
			slog.Error("send response error", "err", err)
		}
	}
}

func spamFilterMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil && strings.Contains(update.Message.Text, "spam") {
			return
		}
		next(ctx, b, update)
	}
}

func extractPayload(text string) string {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func runExternalModule(db db.ZbotDatabase, message *models.Message, modules []ExternalModule) string {

	cmd, err := utils.ParseCommand(message.Text)
	if err != nil {
		slog.Error("parse command error", "err", err)
		return ""
	}

	cmdFile, err := utils.GetCommandFile(cmd, modules)
	if err != nil {
		slog.Error("get command file error", "err", err)
		return ""
	}

	fullPathToBinary, _ := command.LookPathCommand(ModulesPath + cmdFile)

	chatName := ""
	if message.Chat.Type != "private" {
		chatName = message.Chat.Title
	}

	user := user.BuildUser(message.From, db)
	slog.Debug("Running module", "path", fullPathToBinary)
	payload := extractPayload(message.Text)
	response := utils.RunExternalCommand(fullPathToBinary, user.Username, strconv.Itoa(user.Level), chatName, payload)
	return response
}

// messagesProcessing
func messagesProcessing(db db.ZbotDatabase, message *models.Message, chatName string) string {

	private := false
	if message.Chat.Type == "private" && chatName == "" {
		private = true
	}

	//we're going to process only the message starting with ! or ?
	processingMsg := regexp.MustCompilePOSIX(`^[!|?].*`)
	username := strings.ToLower(message.From.Username)

	if !checkIgnoreList(db, username) {
		if processingMsg.MatchString(message.Text) {
			slog.Debug("received message", "user", username, "text", message.Text)
			return cmdProcessing(db, *message, chatName, private)
		}
	} else {
		slog.Debug("user ignored", "user", username)
	}

	return ""
}

// checkIgnoreList check user in the ignore list
// return true if user is on the ignore_list
//
//	false if the flag ignore is disable or the user is not in the list
func checkIgnoreList(db db.ZbotDatabase, username string) bool {

	if Flags.Ignore {
		slog.Debug("Checking user", "user", username)
		return db.UserCheckIgnore(username)
	}

	return false
}

// cmdProcessing process message using commands
func cmdProcessing(db db.ZbotDatabase, msg models.Message, chatName string, private bool) string {

	commandName := command.GetCommandInformation(msg.Text)

	if command.IsCommandDisabled(commandName) {
		slog.Debug("command disabled", "command", commandName)
		return ""
	}

	user := user.BuildUser(msg.From, db)

	if Flags.Level {
		requiredLevel := command.GetMinimumLevel(commandName, levelsConfig)
		if !command.CheckPermission(commandName, user, requiredLevel) {
			return fmt.Sprintf("Your level is not enough < %d", requiredLevel)
		}
	}

	commandsList := &command.CommandsList{
		List: list.New(),
	}

	commandsList.Chain("ping", &command.PingCommand{Db: db}, levelsConfig.Ping)
	commandsList.Chain("version", &command.VersionCommand{
		GitHash:   gitHash,
		Version:   version,
		BuildTime: buildTime,
	}, levelsConfig.Version)
	commandsList.Chain("top", &command.TopCommand{Db: db}, levelsConfig.Top)
	commandsList.Chain("stats", &command.StatsCommand{Db: db}, levelsConfig.Stats)
	commandsList.Chain("last", &command.LastCommand{Db: db}, levelsConfig.Last)
	commandsList.Chain("rand", &command.RandCommand{Db: db}, levelsConfig.Rand)
	commandsList.Chain("who", &command.WhoCommand{Db: db}, levelsConfig.Who)
	commandsList.Chain("find", &command.FindCommand{Db: db}, levelsConfig.Find)
	commandsList.Chain("get", &command.GetCommand{Db: db}, levelsConfig.Get)
	commandsList.Chain("search", &command.SearchCommand{Db: db}, levelsConfig.Search)
	commandsList.Chain("learn", &command.LearnCommand{Db: db}, levelsConfig.Learn)
	commandsList.Chain("append", &command.AppendCommand{Db: db}, levelsConfig.Append)
	commandsList.Chain("forget", &command.ForgetCommand{Db: db}, levelsConfig.Forget)
	commandsList.Chain("level", &command.LevelCommand{Db: db}, levelsConfig.Level)
	commandsList.Chain("lock", &command.LockCommand{Db: db}, levelsConfig.Lock)
	commandsList.Chain("ignore", &command.IgnoreCommand{Db: db}, levelsConfig.Ignore)

	var messageString = msg.Text

	if msg.ReplyToMessage != nil {
		messageString = fmt.Sprintf("%s %s %s", messageString, msg.ReplyToMessage.From.Username, msg.ReplyToMessage.Text)
	}

	outputMsg := commandsList.Run(commandName, messageString, user, chatName, private)

	return outputMsg
}

// SetDisabledCommands setup disabled commands
func SetDisabledCommands(cmdList []string) {
	command.DisabledCommands = cmdList
}

// GetDisabledCommands get disabled zbot commands
func GetDisabledCommands() []string {
	return command.DisabledCommands
}

// SetDisabledLearnChannels set list of channels where learns commands wont be used
func SetDisabledLearnChannels(channelsList []string) {
	command.DisableLearnChannels = channelsList
}

func appendChannel(channels []Channel, chat models.Chat) []Channel {

	for i := range channels {
		if channels[i].ID == chat.ID {
			channels[i].Title = chat.Title
			return channels
		} else if channels[i].ID == 0 && chat.Title == channels[i].Title {
			channels[i].ID = chat.ID
			return channels
		}
	}

	channels = append(channels, Channel{ID: chat.ID, Title: chat.Title})
	return channels
}
