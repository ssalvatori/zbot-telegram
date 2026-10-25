package zbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAdmin(t *testing.T) {
	b := &Bot{
		cfg: Config{
			Admins: []int64{123},
		},
	}

	tests := []struct {
		name   string
		userID int64
		want   bool
	}{
		{"admin user", 123, true},
		{"other user", 456, false},
		{"zero ID", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, b.IsAdmin(tt.userID), "IsAdmin(%d)", tt.userID)
		})
	}
}
