package command

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"

	"strings"

	"github.com/ssalvatori/zbot-telegram-go/user"
)

//LevelCommand definition
type LevelCommand struct {
	Db db.ZbotDatabase
}

//AddUser add level for a given user
func (handler *LevelCommand) AddUser(targetUser string, level int, responsibleUser string) (string, error) {
	return "not ready", nil
}

//DelUser delete level for a given user
func (handler *LevelCommand) DelUser(userToCheck string, user string) (string, error) {
	return "not ready", nil
}

//GetLevel return level for a given user
func (handler *LevelCommand) GetLevel(user string) (string, error) {
	level, err := handler.Db.UserLevel(user)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return fmt.Sprintf("%s level %s", user, level), nil
}

//PaserCommand identify user to check levels
func (handler *LevelCommand) PaserCommand(cmd string, user string) map[string]string {

	res := make(map[string]string)
	parsed := strings.Fields(cmd)

	res["user"] = user
	res["subcommand"] = "get"
	res["level"] = "0"

	if len(parsed) == 4 {
		switch parsed[1] {
		case "add":
			res["subcommand"] = "add"
			res["user"] = parsed[2]
			res["level"] = parsed[3]
		default:

		}
	}

	if len(parsed) == 3 {
		switch parsed[1] {
		case "del":
			res["subcommand"] = "del"
			res["user"] = parsed[2]
			res["level"] = "0"
		default:

		}
	}

	return res
}

//ProcessText run command
func (handler *LevelCommand) ProcessText(text string, user user.User) (string, error) {
	commandPattern := regexp.MustCompile(`^!level(\s|$)(\S*)\s?(\S+)?\s?(\d+)?`)
	var result string
	var err error

	if commandPattern.MatchString(text) {
		parsedCmd := handler.PaserCommand(text, user.Username)
		switch parsedCmd["subcommand"] {
		case "add":
			level, _ := strconv.Atoi(parsedCmd["level"])
			result, err = handler.AddUser(parsedCmd["user"], level, user.Username)
		case "del":
			result, err = handler.DelUser(parsedCmd["user"], user.Username)
		default:
			result, err = handler.GetLevel(parsedCmd["user"])
		}
		if err != nil {
			return "", err
		}

		return result, nil
	}
	return "", errors.New("text doesn't match")
}
