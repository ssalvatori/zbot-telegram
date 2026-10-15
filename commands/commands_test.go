package command

import (
	"container/list"
	"errors"
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/stretchr/testify/assert"
)

// userTest is shared by every command test and must not be mutated.
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

// assertCommandError checks ProcessText's error against a sentinel, any error, or no error.
func assertCommandError(t *testing.T, err error, wantErr error, wantAnyErr bool) {
	t.Helper()
	switch {
	case wantErr != nil:
		assert.ErrorIs(t, err, wantErr)
	case wantAnyErr:
		assert.Error(t, err)
	default:
		assert.NoError(t, err)
	}
}

// disableLearn marks testchat as learn-disabled for the duration of the test.
func disableLearn(t *testing.T) {
	t.Helper()
	DisableLearnChannels = []string{"testchat"}
	t.Cleanup(func() { DisableLearnChannels = nil })
}

type fakeCommand struct {
	output string
	err    error
}

func (f *fakeCommand) ProcessText(text string, user user.User, channelID int, channelName string, private bool) (string, error) {
	return f.output, f.err
}

func TestGetTerms(t *testing.T) {
	tests := []struct {
		name  string
		items []db.Definition
		want  []string
	}{
		{
			name:  "terms",
			items: []db.Definition{{Term: "foo", Meaning: "bar"}, {Term: "foo2", Meaning: "bar2"}, {Term: "foo3", Meaning: "bar3"}},
			want:  []string{"foo", "foo2", "foo3"},
		},
		{name: "empty", items: []db.Definition{}, want: nil},
		{name: "blank term skipped", items: []db.Definition{{Term: "", Meaning: "bar"}}, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getTerms(tt.items))
		})
	}
}

func TestIsCommandDisabled(t *testing.T) {
	DisabledCommands = []string{"learn", "version"}
	t.Cleanup(func() { DisabledCommands = nil })

	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{"disabled", "learn", true},
		{"enabled", "lock", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsCommandDisabled(tt.command))
		})
	}
}

func TestGetCommandInformation(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"no prefix", "asdfasd", ""},
		{"command", "!hola", "hola"},
		{"lowercased", "!Learn", "learn"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetCommandInformation(tt.text))
		})
	}
}

func TestCheckPermission(t *testing.T) {
	tests := []struct {
		name          string
		level         int
		requiredLevel int
		want          bool
	}{
		{"equal level", 10, 10, true},
		{"higher level", 100, 10, true},
		{"lower level", 5, 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := userTest
			u.Level = tt.level
			assert.Equal(t, tt.want, CheckPermission("cmd", u, tt.requiredLevel))
		})
	}
}

func TestGetMinimumLevel(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    int
	}{
		{"defined", "lock", minimumLevels.Lock},
		{"not defined", "hola", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetMinimumLevel(tt.command, minimumLevels))
		})
	}
}

func TestCommandsListChain(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		wantLen  int
	}{
		{"empty", nil, 0},
		{"single", []string{"one"}, 1},
		{"multiple", []string{"one", "two", "three"}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &CommandsList{List: list.New()}
			for _, name := range tt.commands {
				l.Chain(name, &fakeCommand{}, 0, name+" help")
			}
			assert.Equal(t, tt.wantLen, l.List.Len())
		})
	}
}

func TestCommandsListRun(t *testing.T) {
	tests := []struct {
		name     string
		commands []zbotCommand
		want     string
	}{
		{
			name: "first matching command wins",
			commands: []zbotCommand{
				&fakeCommand{err: ErrNextCommand},
				&fakeCommand{output: "Fake 2"},
				&fakeCommand{output: "Fake 3"},
			},
			want: "Fake 2",
		},
		{
			name:     "no command matches",
			commands: []zbotCommand{&fakeCommand{err: ErrNextCommand}},
			want:     "",
		},
		{
			name:     "learn disabled channel",
			commands: []zbotCommand{&fakeCommand{err: ErrLearnDisabledChannel}},
			want:     "Learn modules are been disabled for this channel!",
		},
		{
			name:     "internal error",
			commands: []zbotCommand{&fakeCommand{err: errors.New("some internal error")}},
			want:     "Internal Error, check logs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &CommandsList{List: list.New()}
			for _, cmd := range tt.commands {
				l.Chain("cmd", cmd, 0, "")
			}
			assert.Equal(t, tt.want, l.Run("cmd", "!cmd", userTest, 0, "test", false))
		})
	}
}

func TestSetup(t *testing.T) {
	Setup()
	assert.IsNonDecreasing(t, LearnCommandList)
}

func TestCheckLearnCommandOnChannel(t *testing.T) {
	DisableLearnChannels = []string{"disabled-chat"}
	t.Cleanup(func() { DisableLearnChannels = nil })

	tests := []struct {
		name    string
		channel string
		want    bool
	}{
		{"disabled", "disabled-chat", true},
		{"enabled", "testchat", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, checkLearnCommandOnChannel(tt.channel))
		})
	}
}
