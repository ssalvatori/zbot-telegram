package command

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

// TopCommand definition
type TopCommand struct {
	Db db.ZbotDatabase
}

//SetDb set db connection if the module need it
func (handler *TopCommand) SetDb(db db.ZbotDatabase) {}

// ProcessText run command
func (handler *TopCommand) ProcessText(text string, user user.User, chat string) (string, error) {

	commandPattern := regexp.MustCompile(`^!top(\s+(\d+))?$`)

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		var limit int = 10
		var err error = nil

		if len(term) == 3 && term[2] != "" {
			limit, err = strconv.Atoi(term[2])

			if err != nil {
				log.Error(fmt.Printf("Problem converting %s", term[2]))
				limit = 10
			}
		}

		if limit > 100 {
			limit = 100
		}

		items, err := handler.Db.Top(limit)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("no results"), nil
			}
			log.Error(err)
			return "", fmt.Errorf("Internal error, check logs")
		}
		return fmt.Sprintf(strings.Join(getTerms(items), " ")), nil
	}

	return "", ErrNextCommand

}
