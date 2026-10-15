package zbot

import (
	"testing"
)

func TestIsAdmin(t *testing.T) {
	b := &Bot{
		cfg: Config{
			Admins: []string{"stefano_bikes"},
		},
	}

	tests := []struct {
		username string
		expected bool
	}{
		{"stefano_bikes", true},
		{"other_user", false},
	}

	for _, test := range tests {
		result := b.IsAdmin(test.username)
		if result != test.expected {
			t.Errorf("IsAdmin(%s) = %v; want %v", test.username, result, test.expected)
		}
	}
}
