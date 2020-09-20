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

//RandCommand definition
type RandCommand struct {
	Db db.ZbotDatabase
}

// ProcessText run command
func (handler *RandCommand) ProcessText(text string, user user.User, chat string) (string, error) {

	commandPattern := regexp.MustCompile(`^!rand(\s+(\d+))?$`)

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		var limit int = 1
		var err error = nil

		if len(term) == 3 && term[2] != "" {
			limit, err = strconv.Atoi(term[2])

			if err != nil {
				log.Error(fmt.Printf("Problem converting %s", term[2]))
				limit = 1
			}
		}

		if limit > 100 {
			limit = 100
		}

		items, err := handler.Db.Rand(limit)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return fmt.Sprintf("no results"), nil
			}
			log.Error(err)
			return "", fmt.Errorf("Internal error, check logs")
		}
		var output []string
		for _, item := range items {
			output = append(output, fmt.Sprintf("[%s] - [%s]", item.Term, item.Meaning))

		}

		return strings.Join(output, "\n\n\n"), nil
	}
	return "", ErrNextCommand
}
