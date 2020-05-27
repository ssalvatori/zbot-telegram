package zbot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"container/list"

	log "github.com/sirupsen/logrus"
	command "github.com/ssalvatori/zbot-telegram/commands"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	version   = "dev-master"
	buildTime = time.Now().String()
	gitHash   = "undefined"
	//DatabaseType database backend to be use (mysql or sqlite)
	DatabaseType = ""
	//APIToken Telegram API Token (key:secret Format)
	APIToken = ""
	//ModulesPath Absolute path where the modules are located
	ModulesPath = utils.GetCurrentDirectory() + "/../modules/"
	//Flags zbot configurations
	Flags = ConfigurationFlags{Ignore: false, Level: false}
)

//ConfigurationFlags configurations false means the feature is disabled
type ConfigurationFlags struct {
	Ignore bool
	Level  bool
}

//Db interface to the database
var Db db.ZbotDatabase

var levelsConfig = command.Levels{
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

//Execute run Zbot
func Execute() {
	log.Info("Loading zbot-telegram version [" + version + "] [" + buildTime + "] [" + gitHash + "]")

	log.Info("Database: ", DatabaseType)
	log.Info("Modules: ", ModulesPath)
	log.Info("Configuration Flags Ignore: ", Flags.Ignore)
	log.Info("Configuration Flags Level: ", Flags.Level)

	bot, err := tb.NewBot(tb.Settings{
		Token:  APIToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	err = Db.Init()
	defer Db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if Flags.Ignore {
		go Db.UserCleanIgnore()
	}

	bot.Handle(tb.OnText, func(m *tb.Message) {
		var response = messagesProcessing(Db, m)
		if response != "" {
			bot.Send(m.Chat, response)
		}
	})

	bot.Start()
}

//messagesProcessing
func messagesProcessing(db db.ZbotDatabase, message *tb.Message) string {

	//we're going to process only the message starting with ! or ?
	processingMsg := regexp.MustCompilePOSIX(`^[!|?].*`)
	username := strings.ToLower(message.Sender.Username)

	if !checkIgnoreList(db, username) {
		if processingMsg.MatchString(message.Text) {
			log.Debug(fmt.Sprintf("Received a message from %s with the text: %s", username, message.Text))
			return processing(db, *message)
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

// processing process message using commands
func processing(db db.ZbotDatabase, msg tb.Message) string {

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
	commandsList.Chain("external", &command.ExternalCommand{PathModules: ModulesPath}, levelsConfig.External)

	var messageString = msg.Text

	if msg.ReplyTo != nil {
		messageString = fmt.Sprintf("%s %s %s", messageString, msg.ReplyTo.Sender.Username, msg.ReplyTo.Text)
	}

	outputMsg := commandsList.Run(commandName, messageString, user)

	//	outputMsg := commands.ProcessText(messageString, user)

	return outputMsg
}

//SetDisabledCommands setup disabled commands
func SetDisabledCommands(dataBinaryContent []byte) {
	var c []string
	err := json.Unmarshal(dataBinaryContent, &c)

	if err != nil {
		log.Debug("No disabled commands")
		command.DisabledCommands = []string{}
	}

	command.DisabledCommands = c
	sort.Strings(command.DisabledCommands)
}

//GetDisabledCommands get disabled zbot commands
func GetDisabledCommands() []string {
	return command.DisabledCommands
}
