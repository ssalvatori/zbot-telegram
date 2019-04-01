package command

import (
	"container/list"
	"errors"
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/ssalvatori/zbot-telegram-go/user"
	"github.com/stretchr/testify/assert"
)

// userTest variable will be shared between all the test
var userTest = user.User{
	Username: "ssalvatori",
	Ident:    "stefano",
	Host:     "localhost",
	Level:    100,
}

var minimumLevels = Levels{
	Ignore: 100,
	Lock:   1000,
	Learn:  10,
	Append: 0,
	Forget: 1000,
	Who:    0,
	Top:    0,
	Stats:  0,
}

type FakeCommand struct {
	Next HandlerCommand
}

func (handler *FakeCommand) ProcessText(text string, user user.User) string {
	return "Fake OK"
}

func TestGetTerms(t *testing.T) {
	items := []db.DefinitionItem{
		{Term: "foo", Meaning: "bar"},
		{Term: "foo2", Meaning: "bar2"},
		{Term: "foo3", Meaning: "bar3"},
	}
	assert.Equal(t, []string{"foo", "foo2", "foo3"}, getTerms(items))
	var terms []string
	assert.Equal(t, terms, getTerms([]db.DefinitionItem{}))
}

func TestIsCommandDisable(t *testing.T) {

	DisabledCommands = []string{
		"learn",
		"version",
	}

	assert.True(t, IsCommandDisabled("learn"), "Disable Commands true")

	assert.False(t, IsCommandDisabled("lock"), "Disable Commands false")
}

func TestGetCommandInformation(t *testing.T) {
	assert.Equal(t, "", GetCommandInformation("asdfasd"), "Command not found")
	assert.Equal(t, "hola", GetCommandInformation("!hola"), "Command not found")

	assert.Equal(t, "learn", GetCommandInformation("!Learn"), "Command not found")
}

func TestCheckPermission(t *testing.T) {
	userTest.Level = 10
	assert.True(t, CheckPermission("hola", userTest, 10), "")
	userTest.Level = 5
	assert.False(t, CheckPermission("learn", userTest, 1000), "")
}

func TestGetMinimumLevel(t *testing.T) {

	assert.Equal(t, minimumLevels.Lock, GetMinimumLevel("lock", minimumLevels), "checking lock")

	assert.Equal(t, 0, GetMinimumLevel("hola", minimumLevels), "checking level not defined")
}

type FakeCommand2 struct {
	cmd string
}

func (handler *FakeCommand2) ProcessText(text string, user user.User) (string, error) {
	if text != "!fakecommand2" {
		return "", errors.New("text doesn't match")
	}
	return "Fake 2", nil
}

func (handler *FakeCommand2) SetDB(db db.ZbotDatabase) {

}

var fake1Command = &FakeCommand2{cmd: "fakecommand1"}
var fake2Command = &FakeCommand2{cmd: "fakecommand2"}
var fake3Command = &FakeCommand2{cmd: "fakecommand3"}

var cmdList = &CommandsList{
	List: list.New(),
	Db: &db.MockZbotDatabase{
		Term:    "foo",
		Meaning: "bar",
	},
}

func TestChainAndRun(t *testing.T) {

	cmdList.Chain("fakecommand1", fake1Command, 0)
	cmdList.Chain("fakecommand2", fake2Command, 0)
	cmdList.Chain("fakecommand3", fake2Command, 0)
	assert.Equal(t, 3, cmdList.List.Len(), "Chain add elements in the list")

	assert.Equal(t, "Fake 2", cmdList.Run("fakecommand2", "!fakecommand2", userTest))
	assert.Equal(t, "text doesn't match", cmdList.Run("fakecommand1", "!fakecommand1", userTest))
	assert.Equal(t, "text doesn't match", cmdList.Run("fakecommand2", "!fakecommand5", userTest))

}

/*
func TestSetDisabledCommands(t *testing.T) {

	commands := `["level","ignore"]`
	jsonRaw := json.RawMessage(commands)
	binary, _ := jsonRaw.MarshalJSON()
	SetDisabledCommands(binary)
	disabledCommands := []string{"ignore", "level"}

	assert.Equal(t, disabledCommands, DisabledCommands, "disabled command")

	DisabledCommands = []string(nil)

	commands = ``
	jsonRaw = json.RawMessage(commands)
	binary, _ = jsonRaw.MarshalJSON()
	SetDisabledCommands(binary)
	disabledCommands = []string(nil)

	assert.Equal(t, disabledCommands, DisabledCommands, "no disabled command")
}
*/
