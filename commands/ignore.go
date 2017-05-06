package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Levels struct {
	Ignore   int
	Lock     int
	Append   int
	Learn    int
	Forget   int
	Who	int
	LevelAdd int
	LevelDel int
}
type IgnoreCommand struct {
	Next   HandlerCommand
	Db     db.ZbotDatabase
	Levels Levels
}

const dateFormat string = "02-01-2006 15:04:05 MST" //dd-mm-yyyy hh:ii:ss TZ

func (handler *IgnoreCommand) ProcessText(text string, user user.User) string {
	commandPattern := regexp.MustCompile(`^!ignore\s(\S*)(\s(\S*))?`)
	result := ""

	if commandPattern.MatchString(text) {
		args := commandPattern.FindStringSubmatch(text)

		switch args[1] {
		case "help":
			result = "*!ignore* Options available: \n list (show all user ignored) \n add <username> (ignore a user for 10 minutes)"
			break
		case "list":
			ignoredUsers, err := handler.Db.UserIgnoreList()
			if err != nil {
				log.Error(err)
			}
			result = fmt.Sprintf(strings.Join(getUserIgnored(ignoredUsers), "/n"))
			break
		case "add":
			if user.IsAllow(handler.Levels.Ignore) {
				if strings.ToLower(args[3]) != strings.ToLower(user.Username) {
					err := handler.Db.UserIgnoreInsert(args[3])
					if err != nil {
						log.Error(err)
					}
					result = fmt.Sprintf("User [%s] ignored for 10 minutes", args[3])
				} else {
					result = fmt.Sprintf("You can't ignore youself")
				}
			} else {
				result = fmt.Sprintf("Your level is not enough < %d", handler.Levels.Ignore)
			}
			break
		}
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}

func getUserIgnored(users []db.UserIgnore) []string {
	var usersIgnored []string
	for _, userIgnore := range users {
		if userIgnore.Username != "" {
			since, until := convertDates(userIgnore.Since, userIgnore.Until)
			userString := fmt.Sprintf("[ @%s ] since [%s] until [%s]", userIgnore.Username, since, until)
			usersIgnored = append(usersIgnored, userString)
		}
	}
	return usersIgnored
}

func convertDates(since string, until string) (string, string) {
	time.LoadLocation("UTC")

	i, err := strconv.ParseInt(since, 10, 64)
	sinceFormated := time.Unix(100, 0)
	untilFormated := time.Unix(600, 0)
	if err != nil {
		log.Error("converting ignore time (since)")
	} else {
		sinceFormated = time.Unix(i, 0)
	}

	i, err = strconv.ParseInt(until, 10, 64)
	if err != nil {
		log.Error("converting ignore time (until)")
	} else {
		untilFormated = time.Unix(i, 0)
	}

	return sinceFormated.UTC().Format(dateFormat), untilFormated.UTC().Format(dateFormat)
}
