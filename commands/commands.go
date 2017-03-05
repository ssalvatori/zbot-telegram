package command

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"strconv"
)

type HandlerCommand interface {
	ProcessText(text string, username User) string
}

type User struct {
	Username string
	Ident    string
	Host     string
}

func getTerms(items []db.DefinitionItem) []string {
	var terms []string
	for _, item := range items {
		if item.Term != "" {
			terms = append(terms, item.Term)
		}
	}
	return terms
}

func IsUserAllow(Db db.ZbotDatabase, username string, level int) bool {
	result := false

	userLevel, err := Db.UserLevel(username)
	if err != nil {
		log.Error(err)
		return false
	}
	userLevelInt, err := strconv.Atoi(userLevel)
	if err != nil {
		log.Error(err)
		return false
	}
	if userLevelInt >= level {
		result = true
	}

	return result
}