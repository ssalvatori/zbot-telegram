package zbot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"container/list"

	log "github.com/sirupsen/logrus"
	command "github.com/ssalvatori/zbot-telegram/commands"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/server"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//ExternalModule definition
type ExternalModule struct {
	Key         string
	File        string
	Description string
}

//Channel definition
type Channel struct {
	ID        int64
	Title     string
	AuthToken string
}

//ConfigurationFlags configurations false means the feature is disabled
type ConfigurationFlags struct {
	Ignore bool
	Level  bool
}

//ConfigurationWebhook configuration
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

	//Webhook configuration
	Webhook = ConfigurationWebhook{Enable: false, Port: 11337}

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

//Execute run Zbot
func Execute() {
	log.Info("Loading zbot-telegram version [" + version + "] [" + buildTime + "] [" + gitHash + "]")

	log.Info("Database: ", DatabaseType)
	log.Info("Modules: ", ModulesPath)
	log.Info("Configuration Flags Ignore: ", Flags.Ignore)
	log.Info("Configuration Flags Level: ", Flags.Level)

	command.Setup()

	poller := &tb.LongPoller{Timeout: 10 * time.Second}

	middleware := tb.NewMiddlewarePoller(poller, middleware)

	bot, err := tb.NewBot(tb.Settings{
		Token:       APIToken,
		Poller:      middleware,
		Synchronous: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = Db.Init()
	defer Db.Close()

	if err != nil {
		log.Fatal(err)
	}

	//TODO: Not implemented
	// if Flags.Ignore {
	// 	//go Db.UserCleanupIgnorelist()
	// }

	log.Debug(fmt.Sprintf("Modules to load %+v", ExternalModules))
	botCommands := []tb.Command{}

	//Register extra modules
	for _, module := range ExternalModules {
		var cmdString = fmt.Sprintf("/%s", module.Key)
		log.Debug(fmt.Sprintf("Loading module %s from path %s%s", module.Key, ModulesPath, module.File))

		_, err := command.LookPathCommand(ModulesPath + module.File)

		if err != nil {
			log.Error(fmt.Sprintf("File %s for module [%s] can't be loaded %s", module.File, module.Key, ModulesPath))
			log.Error(err)
			continue
		}

		bot.Handle(cmdString, func(m *tb.Message) {
			response := runExternalModule(Db, m, ExternalModules)
			_, err = bot.Send(m.Chat, response)
			if err != nil {
				log.Error(err)
				log.Error("Could not send the message")
			}
		})
		botCommands = append(botCommands, tb.Command{Text: "/" + module.Key, Description: module.Description})
	}

	log.Debug(fmt.Sprintf("Seting bot commands: %+v", botCommands))
	err = bot.SetCommands(botCommands)
	if err != nil {
		log.Error("Error trying to set commands")
		log.Error(err)
	}

	bot.Handle(tb.OnText, func(m *tb.Message) {
		chatName := ""
		if m.Chat.Type != "private" {
			chatName = m.Chat.Title
		}

		var response = messagesProcessing(Db, m, chatName)
		if response != "" {
			_, err = bot.Send(m.Chat, response)
			if err != nil {
				log.Error("Could not send the message")
				log.Error(err)
			}
		}
	})

	go bot.Start()
	if Webhook.Enable {
		go server.Start(Webhook.Port, bot, Channels)
	}
	select {} // keep running
}

func runExternalModule(db db.ZbotDatabase, message *tb.Message, modules []ExternalModule) string {

	cmd, err := utils.ParseCommand(message.Text)
	if err != nil {
		log.Error(err)
		return ""
	}

	cmdFile, err := utils.GetCommandFile(cmd, modules)
	if err != nil {
		log.Error(err)
		return ""
	}

	fullPathToBinary, _ := command.LookPathCommand(ModulesPath + cmdFile)

	chatName := ""
	if message.Chat.Type != "private" {
		chatName = message.Chat.Title
	}

	user := user.BuildUser(message.Sender, db)
	log.Debug(fmt.Sprintf("Running module %s ", fullPathToBinary))
	response := utils.RunExternalCommand(fullPathToBinary, user.Username, strconv.Itoa(user.Level), chatName, strings.TrimSpace(message.Payload))
	return response
}

//messagesProcessing
func messagesProcessing(db db.ZbotDatabase, message *tb.Message, chatName string) string {

	private := false
	if message.Chat.Type == "private" && chatName == "" {
		private = true
	}

	//we're going to process only the message starting with ! or ?
	processingMsg := regexp.MustCompilePOSIX(`^[!|?].*`)
	username := strings.ToLower(message.Sender.Username)

	if !checkIgnoreList(db, username) {
		if processingMsg.MatchString(message.Text) {
			log.Debug(fmt.Sprintf("Received a message from %s with the text: %s", username, message.Text))
			return cmdProcessing(db, *message, chatName, private)
		}
	} else {
		log.Debug(fmt.Sprintf("User [%s] ignored", username))
	}

	return ""
}

//checkIgnoreList check user in the ignore list
//return true if user is on the ignore_list
//		 false if the flag ignore is disable or the user is not in the list
func checkIgnoreList(db db.ZbotDatabase, username string) bool {

	if Flags.Ignore {
		log.Debug(fmt.Sprintf("Checking user [%s] ", username))
		return db.UserCheckIgnore(username)
	}

	return false
}

//cmdProcessing process message using commands
func cmdProcessing(db db.ZbotDatabase, msg tb.Message, chatName string, private bool) string {

	commandName := command.GetCommandInformation(msg.Text)

	if command.IsCommandDisabled(commandName) {
		log.Debug("Command [", commandName, "] is disabled")
		return ""
	}

	user := user.BuildUser(msg.Sender, db)

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
	// commandsList.Chain("external", &command.ExternalCommand{PathModules: ModulesPath}, levelsConfig.External)

	var messageString = msg.Text

	if msg.ReplyTo != nil {
		messageString = fmt.Sprintf("%s %s %s", messageString, msg.ReplyTo.Sender.Username, msg.ReplyTo.Text)
	}

	outputMsg := commandsList.Run(commandName, messageString, user, chatName, private)

	//	outputMsg := commands.ProcessText(messageString, user)

	return outputMsg
}

//SetDisabledCommands setup disabled commands
func SetDisabledCommands(cmdList []string) {
	command.DisabledCommands = cmdList
}

//GetDisabledCommands get disabled zbot commands
func GetDisabledCommands() []string {
	return command.DisabledCommands
}

//SetDisabledLearnChannels set list of channels where learns commands wont be used
func SetDisabledLearnChannels(channelsList []string) {
	command.DisableLearnChannels = channelsList
}

func appendChannel(channels []Channel, chat tb.Chat) []Channel {

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

func middleware(msg *tb.Update) bool {
	if msg.Message == nil {
		return true
	}

	if strings.Contains(msg.Message.Text, "spam") {
		return false
	}

	if msg.Message.Chat.Type == "group" || msg.Message.Chat.Type == "supergroup" {
		Channels = appendChannel(Channels, *msg.Message.Chat)
	}

	return true
}
