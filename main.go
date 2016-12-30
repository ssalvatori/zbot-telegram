package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/tucnak/telebot"

	"github.com/ssalvatori/zbot-telegram-go/commands"
	"github.com/ssalvatori/zbot-telegram-go/database"
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

	db := &sqlLite{
		file: dbFile,
	}

	err = db.init()
	defer db.close()

	if err != nil {
		log.Fatal(err)
	}

	go db.userCleanIgnore()
	bot.Messages = make(chan telebot.Message, 1000)
	go messagesProcessing(db, bot)

	bot.Start(1 * time.Second)
}

func messagesProcessing(db zbotDatabase, bot *telebot.Bot) {
	output := make(chan string)
	for message := range bot.Messages {

		//we're going to process only the message starting with ! or ?
		processingMsg := regexp.MustCompilePOSIX(`^[!|?].*`)

		//check if the user isn't on the ignore_list
		ignore, err := db.userCheckIgnore(strings.ToLower(message.Sender.Username))
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

func sendResponse(bot *telebot.Bot, db zbotDatabase, msg telebot.Message, output chan string) {
	response := processing(db, msg,output)
	bot.SendMessage(msg.Chat, response, nil)
}

func processing(db zbotDatabase, msg telebot.Message, output chan string) string {
	//TODO: check chain responsibility design pattern

	var outputMsg string

/*	//versionPattern := regexp.MustCompile(`^!version`)
	learnPattern := regexp.MustCompile(`^!learn\s(\S*)\s(.*)`)
	getPattern := regexp.MustCompile(`^\?\s(\S*)`)
	findPattern := regexp.MustCompile(`^!find\s(\S*)`)
	searchPattern := regexp.MustCompile(`^!search\s(\S*)`)
	topPattern := regexp.MustCompile(`^!top`)
	lastPattern := regexp.MustCompile(`^!last`)
	randPattern := regexp.MustCompile(`^!rand`)
	statsPattern := regexp.MustCompile(`^!stats`)
	//pingPattern := regexp.MustCompile(`^!ping`)*/

	commands := &command.PingCommand{
		Next: &command.VersionCommand{
			Next: &command.StatsCommand{
				Db: db,
			},
			Version: version,
		},
	}

	outputMsg = commands.ProcessText(msg.Text)

/*	//Levels
	levelPattern := regexp.MustCompile(`^!level`)
	ignorePattern := regexp.MustCompile(`^!ignore\s(\S*)`)
	ignoreListPattern := regexp.MustCompile(`^!ignorelist`)

	nowDate := time.Now().Format("2006-01-02")
	var author string
	var authorIdent string
	if msg.Sender.Username != "" {
		author = msg.Sender.Username
		authorIdent = strings.ToLower(msg.Sender.FirstName)
	} else {
		author = msg.Sender.FirstName
		authorIdent = strings.ToLower(msg.Sender.FirstName)
	}*/

/*	switch {
	case ignoreListPattern.MatchString(msg.Text):
		results, err := db.userIgnoreList()
		if (err != nil) {
			log.Error(err)
			break
		}
		outputMsg = fmt.Sprintf(strings.Join(getUserIgnored(results), "/n"))
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
		result := learnPattern.FindStringSubmatch(msg.Text)
		if author != "" {
			def := definitionItem{
				term:    result[1],
				meaning: result[2],
				author:  fmt.Sprintf("%s!%s@telegram.bot", author, authorIdent),
				date:    nowDate,
			}
			err := db.set(def)
			if err != nil {
				log.Error(err)
			}
			outputMsg = fmt.Sprintf("[%s] - [%s]", def.term, def.meaning)
		} else {
			spew.Dump(msg.Sender)
			outputMsg = ""
		}
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
		result := findPattern.FindStringSubmatch(msg.Text)
		results, err := db.find(result[1])
		if err != nil {
			log.Error(err)
			break
		}
		outputMsg = fmt.Sprintf("%s", strings.Join(getTerms(results), " "))
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

func getTerms(items []definitionItem) ([]string) {
	var terms []string
	for _,item := range items {
		if item.term != "" {
			terms = append(terms, item.term)
		} else {
			return terms
		}
	}
	return terms
}

func getUserIgnored(users []userIgnore) ([]string) {
	var userIgnored []string
	for _,user := range users {
		if user.username != "" {
			userString := fmt.Sprintf("[ @%s ] since [%s] until [%s]", user.username, user.since, user.until)
			userIgnored = append(userIgnored, userString)
		}
	}
	return userIgnored
}
