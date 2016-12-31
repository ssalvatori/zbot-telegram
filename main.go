package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/tucnak/telebot"

	"github.com/ssalvatori/zbot-telegram-go/commands"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"os"
	"time"
	"regexp"
	"strings"
)

const version string = "1.0"
const dbFile string = "./sample.db"
const levelIgnore int = 500 //level minimum to ignore a user

func main() {

	log.Info("Loading zbot-telegram")
	log.SetLevel(log.DebugLevel)

	bot, err := telebot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	db := &db.SqlLite{
		File: dbFile,
	}

	err = db.Init()
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	go db.UserCleanIgnore()
	bot.Messages = make(chan telebot.Message, 1000)
	go messagesProcessing(db, bot)

	bot.Start(1 * time.Second)
}

func messagesProcessing(db db.ZbotDatabase, bot *telebot.Bot) {
	output := make(chan string)
	for message := range bot.Messages {

		//we're going to process only the message starting with ! or ?
		processingMsg := regexp.MustCompilePOSIX(`^[!|?].*`)

		//check if the user isn't on the ignore_list
		ignore, err := db.UserCheckIgnore(strings.ToLower(message.Sender.Username))
		if err != nil {
			log.Error(err)
		}
		if !ignore {
			if processingMsg.MatchString(message.Text) {
				log.Printf("Received a message from %s with the text: %s\n", message.Sender.Username, message.Text)
				go sendResponse(bot, db, message, output)
			}
		} else {
			log.Debug(fmt.Sprintf("User [%s] ignored", strings.ToLower(message.Sender.Username)))
		}
	}
}

func sendResponse(bot *telebot.Bot, db db.ZbotDatabase, msg telebot.Message, output chan string) {
	response := processing(db, msg,output)
	bot.SendMessage(msg.Chat, response, nil)
}

func processing(db db.ZbotDatabase, msg telebot.Message, output chan string) string {

	user := command.User{}

	if msg.Sender.Username != "" {
		user.Username = msg.Sender.Username
		user.Ident = strings.ToLower(msg.Sender.FirstName)
	} else {
		user.Username = msg.Sender.FirstName
		user.Ident = strings.ToLower(msg.Sender.FirstName)
	}

	// TODO: how to clean this code
	commands := &command.PingCommand{}
	versionCommand := &command.VersionCommand{Version: version}
	statsCommand := &command.StatsCommand{Db: db}
	randCommand := &command.RandCommand{Db: db}
	topCommand := &command.TopCommand{Db: db}
	lastCommand := &command.LastCommand{Db: db}
	getCommand := &command.GetCommand{Db: db}
	findCommand :=  &command.FindCommand{Db: db}
	searchCommand := &command.SearchCommand{Db: db}
	learnCommand := &command.LearnCommand{Db: db}
	levelCommand := &command.LevelCommand{Db: db}
	ignoreCommand := &command.IgnoreCommand{Db: db}


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
	levelCommand.Next = ignoreCommand

	outputMsg := commands.ProcessText(msg.Text, user)

/*	//Levels
	ignorePattern := regexp.MustCompile(`^!ignore\s(\S*)`)
	ignoreListPattern := regexp.MustCompile(`^!ignorelist`)

*/

/*	switch {
	case ignoreListPattern.MatchString(msg.Text):

		break
	case ignorePattern.MatchString(msg.Text):
		result := ignorePattern.FindStringSubmatch(msg.Text)
		level, err := db.userLevel(msg.Sender.Username)
		if err != nil {
			log.Error(err)
			break
		}
		levelInt, _ := strconv.Atoi(level)
		if levelInt >= levelIgnore {
			if strings.ToLower(result[1]) != strings.ToLower(msg.Sender.Username) {
				err := db.userIgnoreInsert(result[1])
				if err != nil {
					log.Error(err)
					break
				}
				outputMsg = fmt.Sprintf("User [%s] ignored for 10 minutes", result[1])
			} else {
				outputMsg = fmt.Sprintf("You can't ignore youself")
			}
		} else {
			outputMsg = fmt.Sprintf("level not enough (minimum %d yours %s)", levelIgnore, level)
		}
		break
	case pingPattern.MatchString(msg.Text):
		outputMsg = fmt.Sprintf("pong!!")
		break
	case learnPattern.MatchString(msg.Text):

		break
	case getPattern.MatchString(msg.Text):
		result := getPattern.FindStringSubmatch(msg.Text)
		definition, err := db.get(strings.ToLower(result[1]))
		if err != nil {
			log.Error(err)
			break
		}
		if definition.term != "" {
			outputMsg = fmt.Sprintf("[%s] - [%s]", definition.term, definition.meaning)
		} else {
			outputMsg = fmt.Sprintf("[%s] Not found!", result[1])
		}
		break
	case findPattern.MatchString(msg.Text):

		break
	case searchPattern.MatchString(msg.Text):
		result := searchPattern.FindStringSubmatch(msg.Text)
		results, err := db.search(result[1])
		if err != nil {
			log.Error(err)
			break
		}
		outputMsg = fmt.Sprintf("%s", strings.Join(getTerms(results), " "))
		break
	case topPattern.MatchString(msg.Text):
		items, err := db.top()
		if err != nil {
			log.Error(err)
		}
		outputMsg = fmt.Sprintf(strings.Join(getTerms(items), " "))
		break
	case lastPattern.MatchString(msg.Text):
		lastItem, err := db.last()
		if err != nil {
			log.Error(err)
			break
		}
		outputMsg = fmt.Sprintf("[%s] - [%s]", lastItem.term, lastItem.meaning)
		break
	case versionPattern.MatchString(msg.Text):
		outputMsg = fmt.Sprintf("zbot golang version %s", version)
		break
	case randPattern.MatchString(msg.Text):
		randItem, err := db.rand()
		if err != nil {
			log.Error(err)
			break
		}
		outputMsg = fmt.Sprintf("[%s] - [%s]", randItem.term, randItem.meaning)
		break
	case statsPattern.MatchString(msg.Text):
		statTotal, err := db.statistics()
		if err != nil {
			log.Error(err)
			break
		}
		outputMsg = fmt.Sprintf("Count: %s", statTotal)
		break
	case levelPattern.MatchString(msg.Text):
		level, err := db.userLevel(msg.Sender.Username)
		if err != nil {
			log.Error(err)
			break
		}
		outputMsg = fmt.Sprintf("%s level %s", msg.Sender.Username, level)
		break
	default:
		outputMsg = ""
		break
	}*/

	return outputMsg
}