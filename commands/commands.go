package command

import (
	"errors"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"container/list"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
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
	Rand     int
	Find     int
	Get      int
	Search   int
	External int
	Level    int
}

var (
	//DisabledCommands list with commands that have to be ignored
	DisabledCommands []string
	//ErrNextCommand the command couldn't perform any actions.
	ErrNextCommand = errors.New("no action in command")
	//ErrInternalError internal error
	ErrInternalError = errors.New("internal zbot error")
	//ErrLearnDisabledChannel learn commands are disabled for a channel
	ErrLearnDisabledChannel = errors.New("Learn commands disabled for this channel")
	//DisableLearnChannels list of channels where learn modules won't be working
	DisableLearnChannels []string
	//LearnCommandList list of modules part of learn command
	LearnCommandList = []string{
		"learn",
		"find",
		"who",
		"get",
		"last",
		"lock",
		"rand",
		"stats",
		"top",
		"forget",
		"append",
		"last",
		"search",
	}
)

// CommandsList list of commandElement
type CommandsList struct {
	List *list.List
}

type commandElement struct {
	requiredLevel int
	command       zbotCommand
	cmdString     string
}

//Setup initail setup
func Setup() {
	sort.Strings(LearnCommandList)
}

// Chain add a command and the required level to use it to the list of command
func (cmdList *CommandsList) Chain(cmdDefinition string, cmd interface{}, level int) *CommandsList {

	newCommand := &commandElement{
		command:       cmd.(zbotCommand),
		requiredLevel: level,
		cmdString:     cmdDefinition,
	}

	cmdList.List.PushBack(newCommand)
	return cmdList
}

// Run commands against a msg for a given user
func (cmdList *CommandsList) Run(cmd string, msg string, user user.User, chat string, private bool) string {
	var output string
	var err error

	for e := cmdList.List.Front(); e != nil; e = e.Next() {
		output, err = e.Value.(*commandElement).command.(zbotCommand).ProcessText(msg, user, chat, private)
		if err != nil {
			if errors.Is(err, ErrNextCommand) {
				continue
			}
			if errors.Is(err, ErrLearnDisabledChannel) {
				return "Learn modules are been disabled for this channel!"
			}
			return "Internal Error, check logs"
		}
		return output
	}

	return output
}

// zbotCommand interface to be implemented by each command
type zbotCommand interface {
	ProcessText(text string, username user.User, chat string, private bool) (string, error)
}

//getTerms transform Definition array into an array of term
func getTerms(items []db.Definition) []string {
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

//CheckPermission validate if the user has enought permisison to use a command (each command has a requiredLevel)
func CheckPermission(command string, user user.User, requiredLevel int) bool {
	log.Debug("Checking permission for [", command, "] and user ", user.Username)

	return user.Level >= requiredLevel
}

// IsCommandDisabled check if a command is in the disable list
func IsCommandDisabled(commandName string) bool {
	log.Debug("Checking if [", commandName, "] is disabled")
	//TODO BUG check DisabledCommands before check the array
	return utils.InArray(commandName, DisabledCommands)
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

//checkLearnCommandOnChannel check if channel is in list of channel where learn commands should be disabled
func checkLearnCommandOnChannel(channel string) bool {
	log.Debug("Checking if channel [", channel, "] is the list of ZBOT_CONFIG_DISABLE_LEARN_CHANNELS")

	i := sort.SearchStrings(DisableLearnChannels, channel)
	if i < len(DisableLearnChannels) && DisableLearnChannels[i] == channel {
		return true
	}

	return false
}
