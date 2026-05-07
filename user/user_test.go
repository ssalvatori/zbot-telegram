package user

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/stretchr/testify/assert"
)

func TestBuildUser(t *testing.T) {
	sender := models.User{}
	var newUser User

	mockDatabase := &db.ZbotDatabaseMock{
		Level: "10",
	}

	sender = models.User{
		FirstName: "stefano",
		Username:  "ssalvatori",
	}

	newUser = User{
		Level:    10,
		Username: "ssalvatori",
		Ident:    "stefano",
	}

	assert.Equal(t, newUser, BuildUser(&sender, mockDatabase), "creating with username")

	sender = models.User{
		FirstName: "stefano",
		Username:  "",
	}
	newUser = User{
		Level:    10,
		Username: "stefano",
		Ident:    "stefano",
	}

	assert.Equal(t, newUser, BuildUser(&sender, mockDatabase), "creating without username")

}

func TestGetUserLevel(t *testing.T) {
	userTest := User{Username: "ssalvatori"}

	mockDatabase := &db.ZbotDatabaseMock{
		Level: "10",
	}

	assert.Equal(t, 10, GetUserLevel(mockDatabase, userTest.Username), "Getting user level")

	mockDatabase = &db.ZbotDatabaseMock{
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
