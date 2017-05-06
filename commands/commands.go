package command

import (
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
)
type Levels struct {
	Ignore   int
	Lock     int
	Append   int
	Learn    int
	Forget   int
	Who      int
	LevelAdd int
	LevelDel int
	Top      int
	Stats    int
}

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
