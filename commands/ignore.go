package command

import (
	"regexp"
	"fmt"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"strings"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type Levels struct {
	Ignore int
}
type IgnoreCommand struct {
	Next HandlerCommand
	Db db.ZbotDatabase
	Levels Levels
}

func (handler *IgnoreCommand) ProcessText(text string, user User) string{
	commandPattern := regexp.MustCompile(`^!ignore\s(\S*)(\s(\S*))?`)
	result := ""


	if(commandPattern.MatchString(text)) {
		args := commandPattern.FindStringSubmatch(text)

		switch args[1] {
		case "help":
			result = "*!ignore* Options available: \n list (show all user ignored) add <username> (ignore a user for 10 minutes)"
			break
		case "list":
			ignoredUsers, err := handler.Db.UserIgnoreList()
			if (err != nil) {
				log.Error(err)
			}
			result = fmt.Sprintf(strings.Join(getUserIgnored(ignoredUsers), "/n"))
			break
		case "add":
			level, err := handler.Db.UserLevel(user.Username)
			if err != nil {
				log.Error(err)
			}
			levelInt, _ := strconv.Atoi(level)
			if levelInt >= handler.Levels.Ignore {
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
				result = fmt.Sprintf("level not enough (minimum %d yours %s)", handler.Levels.Ignore, level)
			}
			break
		}
	} else {
		if (handler.Next != nil) {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}

func getUserIgnored(users []db.UserIgnore) ([]string) {
	var userIgnored []string
	for _,user := range users {
		if user.Username != "" {
			userString := fmt.Sprintf("[ @%s ] since [%s] until [%s]", user.Username, user.Since, user.Until)
			userIgnored = append(userIgnored, userString)
		}
	}
	return userIgnored
}

