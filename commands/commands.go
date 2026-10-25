package command

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"

	"container/list"

	"log/slog"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
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
	Help     int
}

var (
	//DisabledCommands list with commands that have to be ignored
	DisabledCommands []string
	//ErrNextCommand the command couldn't perform any actions.
	ErrNextCommand = errors.New("no action in command")
	//ErrInternalError internal error
	ErrInternalError = errors.New("internal zbot error")
	//ErrLearnDisabledChannel learn commands are disabled for a channel
	ErrLearnDisabledChannel = errors.New("learn commands disabled for this channel")
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
// CommandsList represents a list of commands with their required levels.
type CommandsList struct {
	List *list.List
}

type CommandElement struct {
	requiredLevel int
	command       zbotCommand
	CmdString     string
}

// Setup initail setup
func Setup() {
	sort.Strings(LearnCommandList)
}

func NewCommandElement(cmdString string, level int, command zbotCommand) *CommandElement {
	return &CommandElement{
		command:       command,
		requiredLevel: level,
		CmdString:     cmdString,
	}
}

// ProcessText runs the wrapped command against msg.
func (c *CommandElement) ProcessCommand(msg string, u user.User, channelID int, channelName string, private bool) (string, error) {
	return c.command.ProcessText(msg, u, channelID, channelName, private, true)
}

// ShowHelp returns the wrapped command's help message.
func (c *CommandElement) ShowHelp(args ...string) string { return c.command.ShowHelp(args...) }

// ShowDescription returns the wrapped command's description.
func (c *CommandElement) ShowDescription() string { return c.command.ShowDescription() }

// Chain add a command and the required level to use it to the list of command
func (cmdList *CommandsList) Chain(cmdDefinition string, cmd any, level int) *CommandsList {

	newCommand := &CommandElement{
		command:       cmd.(zbotCommand),
		requiredLevel: level,
		CmdString:     cmdDefinition,
	}

	cmdList.List.PushBack(newCommand)
	return cmdList
}

// Run commands against a msg for a given user
func (cmdList *CommandsList) Run(cmd string, msg string, user user.User, channelId int, channelName string, private bool) string {
	var output string
	var err error

	slog.Debug("Running command", "cmd", cmd, "channelId", channelId, "channelName", channelName, "private", private)

	for e := cmdList.List.Front(); e != nil; e = e.Next() {
		output, err = e.Value.(*CommandElement).command.ProcessText(msg, user, channelId, channelName, private, false)
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
	ProcessText(text string, username user.User, channelID int, channelName string, private bool, isCommand bool) (string, error)
	ShowHelp(args ...string) string
	ShowDescription() string
}

// getTerms transform Definition array into an array of term
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

// CheckPermission validate if the user has enought permisison to use a command (each command has a requiredLevel)
func CheckPermission(command string, user user.User, requiredLevel int) bool {
	slog.Debug("Checking permission", "command", command, "user", user.Username)

	return user.Level >= requiredLevel
}

// IsCommandDisabled check if a command is in the disable list
func IsCommandDisabled(commandName string) bool {
	slog.Debug(fmt.Sprintf("Checking if command is disabled [%v]", commandName))
	return slices.Contains(DisabledCommands, commandName)
}

// GetMinimumLevel get the minimum level required for a git command, if it is not defined return 0
func GetMinimumLevel(commandName string, minimumLevels Levels) int {
	slog.Debug("Getting minimum level", "command", commandName)

	field, ok := reflect.TypeOf(&minimumLevels).Elem().FieldByName(cases.Title(language.English).String(commandName))

	if !ok {
		return 0
	}

	r := reflect.ValueOf(&minimumLevels)
	f := reflect.Indirect(r).FieldByName(field.Name)

	return int(f.Int())
}

// checkLearnCommandOnChannel check if channel is in list of channel where learn commands should be disabled
func checkLearnCommandOnChannel(channel string) bool {

	slog.Debug(fmt.Sprintf("Checking channel [%v] is in the disabled learn list", channel))
	slog.Debug(fmt.Sprintf("Disable learn channels: %v", DisableLearnChannels))

	return slices.Contains(DisableLearnChannels, channel)
}
