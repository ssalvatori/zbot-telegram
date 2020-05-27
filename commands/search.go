package command

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

//SearchCommand definition
type SearchCommand struct {
	Db db.ZbotDatabase
}

//ProcessText Run module
func (handler *SearchCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!search\s(\S*)`)

	if commandPattern.MatchString(text) {
		term := commandPattern.FindStringSubmatch(text)
		results, err := handler.Db.Search(term[1])
		if err != nil {
			log.Error(err)
			return "", err
		}
		return fmt.Sprintf("%s", strings.Join(getTerms(results), " ")), nil
	}
	return "", errors.New("text doesn't match")
}
