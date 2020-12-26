package command

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
)

//IgnoreCommand definition
type IgnoreCommand struct {
	Db db.ZbotDatabase
}

const dateFormat string = "02-01-2006 15:04:05 MST" //dd-mm-yyyy hh:ii:ss TZ

//ProcessText run command
func (handler *IgnoreCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!ignore\s(\S*)(\s(\S*))?`)
	result := ""

	if commandPattern.MatchString(text) {
		args := commandPattern.FindStringSubmatch(text)

		switch args[1] {
		case "help":
			result = "*!ignore* Options available: \n list (show all user ignored) \n add <username> (ignore a user for 10 minutes)"
		case "list":
			ignoredUsers, err := handler.Db.UserIgnoreList()
			if err != nil {
				log.Error(err)
				return "", err
			}
			result = fmt.Sprintf(strings.Join(getUserIgnored(ignoredUsers), "/n"))
		case "add":
			if strings.ToLower(args[3]) != strings.ToLower(user.Username) {
				err := handler.Db.UserIgnoreInsert(args[3])
				if err != nil {
					log.Error(err)
					return "", err
				}
				result = fmt.Sprintf("User [%s] ignored for 10 minutes", args[3])
			} else {
				result = "You can't ignore yourself"
			}
		}
		return result, nil
	}
	return "", ErrNextCommand
}

func getUserIgnored(users []db.UserIgnore) []string {
	var usersIgnored []string
	for _, userIgnore := range users {
		if userIgnore.Username != "" {
			// since, until := convertDates(userIgnore.Since, userIgnore.Until)
			userString := fmt.Sprintf("[ @%s ] since [%v] until [%v]", userIgnore.Username, utils.ConvertToDateToUTC(userIgnore.CreatedAt), utils.ConvertToDateToUTC(userIgnore.ValidUntil))
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
