package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InArray(t *testing.T) {
	tests := []struct {
		lookFor        string
		ArrayOfStrings []string
		expected       bool
	}{
		{"abc", []string{"def", "abc", "ghi"}, true},
		{"a1c", []string{"def", "abc", "ghi"}, false},
		{"abc", []string{}, false},
		{"abc", []string{"abc", "abc", "ghi"}, true},
		{"abc", []string{"def", "aXc", "abc"}, true},
	}

	for _, test := range tests {
		got := InArray(test.lookFor, test.ArrayOfStrings)
		assert.Equal(t, test.expected, got, "Testing in array")
	}
}
