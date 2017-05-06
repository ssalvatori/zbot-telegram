package user

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
	"github.com/tucnak/telebot"
)

func TestBuildUser(t *testing.T) {
	sender := telebot.User{}
	newUser := User{}

	mockDatabase := &db.MockZbotDatabase{
		Level: "10",
	}

	sender = telebot.User{
		FirstName: "stefano",
		Username:  "ssalvatori",
	}

	newUser = User{
		Level:    10,
		Username: "ssalvatori",
		Ident:    "stefano",
	}

	assert.Equal(t, newUser, BuildUser(sender, mockDatabase), "creating with username")

	sender = telebot.User{
		FirstName: "stefano",
		Username:  "",
	}
	newUser = User{
		Level:    10,
		Username: "stefano",
		Ident:    "stefano",
	}

	assert.Equal(t, newUser, BuildUser(sender, mockDatabase), "creating without username")

}

func TestGetUserLevel(t *testing.T) {
	userTest := User{Username: "ssalvatori"}

	mockDatabase := &db.MockZbotDatabase{
		Level: "10",
	}

	assert.Equal(t, 10, GetUserLevel(mockDatabase, userTest.Username), "Getting user level")

	mockDatabase = &db.MockZbotDatabase{
		Level: "10",
		Error: true,
	}

	assert.Equal(t, 0, GetUserLevel(mockDatabase, userTest.Username), "Getting user level")

}

func TestIsAllow(t *testing.T) {
	userTest := User{Level: 100}
	assert.True(t, userTest.IsAllow(10), "User is allowed")
	assert.True(t, userTest.IsAllow(100), "User is allowed")
	assert.False(t, userTest.IsAllow(200), "User is allowed")
}
