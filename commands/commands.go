package command

import (
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)

type HandlerCommand interface {
	ProcessText(text string, username user.User) string
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
