package command

import (
	"regexp"
	"fmt"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"strings"
	log "github.com/Sirupsen/logrus"
)

type IgnoreCommand struct {
	Next HandlerCommand
	Db db.ZbotDatabase
}

func (handler *IgnoreCommand) ProcessText(text string, user User) string{

	commandPattern := regexp.MustCompile(`^!ignore\s(\S*)(\S*)`)
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

