package zbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAdmin(t *testing.T) {
	b := &Bot{
		cfg: Config{
			Admins: []string{"user1"},
		},
	}

	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{"admin user", "user1", true},
		{"other user", "other_user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, b.IsAdmin(tt.username), "IsAdmin(%q)", tt.username)
		})
	}
}
