package command

import (
	"reflect"
	"testing"

	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/stretchr/testify/assert"
)

var levelCommand = LevelCommand{}

func TestLevelCommandOK(t *testing.T) {

	levelCommand.Db = &db.ZbotDatabaseMock{
		Level: "1000",
	}
	result, _ := levelCommand.ProcessText("!level", userTest, "testchat")
	assert.Equal(t, "ssalvatori level 1000", result, "Get Level from the same user")
}

func TestProcessText(t *testing.T) {

	tests := []struct {
		name string
		cmd  string
		user user.User
		want string
	}{
		{"same user", "!level", user.User{Username: "ssalvatori"}, "ssalvatori level 1000"},
		{"add other user", "!level add rigo 10", user.User{Username: "ssalvatori"}, "not ready"},
		{"del other", "!level del rigo", user.User{Username: "ssalvatori"}, "not ready"},
		{"del other", "!level2 del rigo", user.User{Username: "ssalvatori"}, ""},
	}

	levelCommand.Db = &db.ZbotDatabaseMock{
		Level: "1000",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := levelCommand.ProcessText(tt.cmd, tt.user, "testchat"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LevelCommand.ProcessText() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPaserCommand(t *testing.T) {

	tests := []struct {
		name string
		cmd  string
		user string
		want map[string]string
	}{
		{"same user", "!level", "ssalvatori", map[string]string{"subcommand": "get", "user": "ssalvatori", "level": "0"}},
		{"add other user", "!level add rigo 10", "ssalvatori", map[string]string{"subcommand": "add", "user": "rigo", "level": "10"}},
		{"del other", "!level del rigo", "ssalvatori", map[string]string{"subcommand": "del", "user": "rigo", "level": "0"}},
	}

	levelCommand.Db = &db.ZbotDatabaseMock{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := levelCommand.PaserCommand(tt.cmd, tt.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LevelCommand.PaserCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelCommandError(t *testing.T) {

	levelCommand.Db = &db.ZbotDatabaseMock{
		RandDef: []db.Definition{db.Definition{Term: "foo", Meaning: "bar"}},
		Error:   true,
	}

	_, err := levelCommand.ProcessText("!level", userTest, "testchat")
	// assert.Equal(t, "Internal error", err.Error(), "Db error")
	assert.Error(t, err, "Internal Error")
}
