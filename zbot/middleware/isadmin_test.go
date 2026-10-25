package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v4"
)

func TestAdminOnly(t *testing.T) {
	const adminID int64 = 123
	isAdmin := func(userID int64) bool { return userID == adminID }

	tests := []struct {
		name       string
		sender     *tele.User
		wantCalled bool
	}{
		{name: "admin without username passes", sender: &tele.User{ID: adminID}, wantCalled: true},
		{name: "non-admin blocked", sender: &tele.User{ID: 456, Username: "admin"}, wantCalled: false},
		{name: "zero ID blocked", sender: &tele.User{Username: "admin"}, wantCalled: false},
		{name: "no sender blocked", sender: nil, wantCalled: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			next := func(c tele.Context) error {
				called = true
				return nil
			}

			var msg *tele.Message
			if tt.sender != nil {
				msg = &tele.Message{Sender: tt.sender}
			}
			ctx := (&tele.Bot{}).NewContext(tele.Update{Message: msg})

			err := AdminOnly(isAdmin)(next)(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCalled, called)
		})
	}
}
