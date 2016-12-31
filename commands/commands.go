package command

import "github.com/ssalvatori/zbot-telegram-go/db"

type HandlerCommand interface {
	ProcessText(text string) string
}

func getTerms(items []db.DefinitionItem) ([]string) {
	var terms []string
	for _,item := range items {
		if item.Term != "" {
			terms = append(terms, item.Term)
		} else {
			return terms
		}
	}
	return terms
}