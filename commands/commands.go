package command

import (
	"reflect"
	"regexp"
	"strings"

	"container/list"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"github.com/ssalvatori/zbot-telegram-go/utils"
)

// Levels command definition
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
	Version  int
	Ping     int
	Last     int
}

var (
	DisabledCommands []string
)

// CommandsList list of commandElement
type CommandsList struct {
	List *list.List
	Db   db.ZbotDatabase
}

type commandElement struct {
	requiredLevel int
	command       zbotCommand
	cmdString     string
}

// Chain add a command and the required level to use it to the list of command
func (cmdList *CommandsList) Chain(cmdDefinition string, cmd interface{}, level int) *CommandsList {

	cmd.(zbotCommand).SetDB(cmdList.Db)

	newCommand := &commandElement{
		command:       cmd.(zbotCommand),
		requiredLevel: level,
		cmdString:     cmdDefinition,
	}

	cmdList.List.PushBack(newCommand)
	return cmdList
}

// Run commands against a msg for a given user
func (cmdList *CommandsList) Run(cmd string, msg string, user user.User) string {
	output := ""

	for e := cmdList.List.Front(); e != nil; e = e.Next() {
		output, err := e.Value.(commandElement).command.(zbotCommand).ProcessText(msg, user)
		if err != nil {
			output = err.Error()
		}
		return output
	}

	return output
}

// zbotCommand interface to be implemented by each command
type zbotCommand interface {
	SetDB(db db.ZbotDatabase)
	ProcessText(text string, username user.User) (string, error)
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

// SetDisabledCommands get the disabled commands from binary json
/*
func SetDisabledCommands(dataBinaryContent []byte) {
	var c []string
	err := json.Unmarshal(dataBinaryContent, &c)

	if err != nil {
		log.Debug("No disabled commands")
		return
	}

	DisabledCommands = c
	sort.Strings(DisabledCommands)
}
*/

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
	}

	return false

}

// IsCommandDisabled check if a command is in the disable list
func IsCommandDisabled(commandName string) bool {
	log.Debug("Checking if [", commandName, "] is disabled")
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
