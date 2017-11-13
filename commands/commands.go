package command

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"github.com/ssalvatori/zbot-telegram-go/utils"
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

var (
	DisabledCommands []string
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

// GetDisabledCommands Reading file to disable son modules
func GetDisabledCommands(disableCommandFile string) {
	log.Debug("Reading file ", disableCommandFile)
	raw, err := ioutil.ReadFile(disableCommandFile)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	var c []string
	json.Unmarshal(raw, &c)
	DisabledCommands = c
	sort.Strings(DisabledCommands)
}

// GetCommandInformation this will parse the text and return the commands and the level minimum to use it or 0 when is
// not defined a level
func GetCommandInformation(text string) string {
	commandPattern := regexp.MustCompile(`^!(\S*)\s*.*`)
	commandName := ""
	if commandPattern.MatchString(text) {
		cmd := commandPattern.FindStringSubmatch(text)
		commandName = strings.ToLower(cmd[1])
	}
	return commandName
}

func CheckPermission(command string, user user.User, requiredLevel int) bool {
	log.Debug("Checking permission for [", command, "] and user ", user.Username)

	if user.Level >= requiredLevel {
		return true
	} else {
		return false
	}
}

// IsCommandDisabled check if a command is in the disable list
func IsCommandDisabled(commandName string) bool {
	log.Debug("Checking isCommandDisabled: [", commandName, "] is disabled")
	//TODO BUG check DisabledCommands before check the array
	if utils.InArray(commandName, DisabledCommands) {
		return true
	}
	return false
}

// GetMinimumLevel get the minimum level required for a git command, if it is not defined return 0
func GetMinimumLevel(commandName string, minimumLevels Levels) int {
	log.Debug("Getting mininum level for ", commandName)

	field, ok := reflect.TypeOf(&minimumLevels).Elem().FieldByName(strings.Title(commandName))

	if !ok {
		return 0
	}

	r := reflect.ValueOf(&minimumLevels)
	f := reflect.Indirect(r).FieldByName(field.Name)

	return int(f.Int())
}
