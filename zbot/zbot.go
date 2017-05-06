package zbot

import (
	"regexp"

	"github.com/ssalvatori/zbot-telegram-go/commands"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"github.com/tucnak/telebot"

	"fmt"
	"strings"

	"os"
	"time"

	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"

	"sort"

	"github.com/ssalvatori/zbot-telegram-go/utils"
)

var (
	Version          = "dev-master"
	BuildTime        = time.Now().String()
	GitHash          = "undefined"
	Database         = "./sample.db"
	ApiToken         = ""
	ModulesPath      = getCurrentDirectory() + "/../modules/"
	DisabledCommands []string
)

var levelsConfig = command.Levels{
	Ignore: 100,
	Lock:   1000,
	Learn:  0,
	Append: 0,
	Forget: 1000,
	Who: 0,
}

// Reading file to disable son modules
func DisableCommands(disableCommandFile string) {
	log.Debug("Reading file ", disableCommandFile)
	raw, err := ioutil.ReadFile(disableCommandFile)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	var c []string
	json.Unmarshal(raw, &c)
	DisabledCommands = c
	sort.Strings(DisabledCommands)
}

func Execute() {
	log.Info("Loading zbot-telegram version [" + Version + "] [" + BuildTime + "] [" + GitHash + "]")

	log.Info("Database: [" + Database + "] Modules: [" + ModulesPath + "]")

	bot, err := telebot.NewBot(ApiToken)
	if err != nil {
		log.Fatal(err)
	}

	database := &db.SqlLite{
		File: Database,
	}

	err = database.Init()
	defer database.Close()

	if err != nil {
		log.Fatal(err)
	}

	go database.UserCleanIgnore()
	bot.Messages = make(chan telebot.Message, 1000)
	go messagesProcessing(database, bot)

	bot.Start(1 * time.Second)
}

func messagesProcessing(db db.ZbotDatabase, bot *telebot.Bot) {
	output := make(chan string)
	for message := range bot.Messages {

		//we're going to process only the message starting with ! or ?
		processingMsg := regexp.MustCompilePOSIX(`^[!|?].*`)

		//check if the user isn't on the ignore_list
		log.Debug(fmt.Sprintf("Checking user [%s] ", strings.ToLower(message.Sender.Username)))
		ignore, err := db.UserCheckIgnore(strings.ToLower(message.Sender.Username))
		if err != nil {
			log.Error(err)
			ignore = true
		}
		if !ignore {
			if processingMsg.MatchString(message.Text) {
				log.Debug(fmt.Sprintf("Received a message from %s with the text: %s", message.Sender.Username, message.Text))
				go sendResponse(bot, db, message, output)
			}
		} else {
			log.Debug(fmt.Sprintf("User [%s] ignored", strings.ToLower(message.Sender.Username)))
		}
	}
}

func sendResponse(bot *telebot.Bot, db db.ZbotDatabase, msg telebot.Message, output chan string) {
	response := processing(db, msg)
	bot.SendMessage(msg.Chat, response, nil)
}

func isCommandDisable(msg telebot.Message) bool {
	commandPattern := regexp.MustCompile(`^!(\S*)\s*.*`)
	if commandPattern.MatchString(msg.Text) {
		cmd := commandPattern.FindStringSubmatch(msg.Text)
		log.Debug("Searching for ", cmd[1])
		log.Debug(DisabledCommands)
		if utils.InArray(cmd[1], DisabledCommands) {
			return true
		}
		return false
	}
	return false
}

func processing(db db.ZbotDatabase, msg telebot.Message) string {

	if isCommandDisable(msg) {
		log.Debug("Command disable")
		return ""
	}

	user := user.BuildUser(msg.Sender, db)

	// TODO: how to clean this code
	commands := &command.PingCommand{}
	versionCommand := &command.VersionCommand{Version: Version, BuildTime: BuildTime}
	statsCommand := &command.StatsCommand{Db: db, Levels: levelsConfig}
	randCommand := &command.RandCommand{Db: db, Levels: levelsConfig}
	topCommand := &command.TopCommand{Db: db, Levels: levelsConfig}
	lastCommand := &command.LastCommand{Db: db, Levels: levelsConfig}
	getCommand := &command.GetCommand{Db: db, Levels: levelsConfig}
	findCommand := &command.FindCommand{Db: db, Levels: levelsConfig}
	searchCommand := &command.SearchCommand{Db: db, Levels: levelsConfig}
	learnCommand := &command.LearnCommand{Db: db, Levels: levelsConfig}
	levelCommand := &command.LevelCommand{Db: db, Levels: levelsConfig}
	ignoreCommand := &command.IgnoreCommand{Db: db, Levels: levelsConfig}
	lockCommand := &command.LockCommand{Db: db, Levels: levelsConfig}
	appendCommand := &command.AppendCommand{Db: db, Levels: levelsConfig}
	whoCommand := &command.WhoCommand{Db: db, Levels: levelsConfig}
	forgetCommand := &command.ForgetCommand{Db: db, Levels: levelsConfig}

	/*
		TODO: check error handler
		!level add <username>
		!level del <username>
	*/

	externalCommand := &command.ExternalCommand{
		PathModules: ModulesPath,
	}

	commands.Next = versionCommand
	versionCommand.Next = statsCommand
	statsCommand.Next = randCommand
	randCommand.Next = topCommand
	topCommand.Next = lastCommand
	lastCommand.Next = getCommand
	getCommand.Next = findCommand
	findCommand.Next = searchCommand
	searchCommand.Next = learnCommand
	learnCommand.Next = levelCommand
	levelCommand.Next = lockCommand
	lockCommand.Next = appendCommand
	appendCommand.Next = whoCommand
	whoCommand.Next = forgetCommand
	forgetCommand.Next = ignoreCommand
	ignoreCommand.Next = externalCommand

	outputMsg := commands.ProcessText(msg.Text, user)

	return outputMsg
}

func getCurrentDirectory() string {
	ex, err := os.Getwd()
	if err != nil {
		log.Panic(err.Error())
		return ""
	}
	return ex
}
