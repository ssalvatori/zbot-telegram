package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInArray(t *testing.T) {
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

func TestGetCurrentDirectory(t *testing.T) {
	currentPath := os.Getenv("PWD")

	assert.Equal(t, GetCurrentDirectory(), currentPath, "")

}

func TestConvertToDate(t *testing.T) {
	date := ConvertToDateToUTC(1604225446)

	assert.Equal(t, "2020-11-01 11:10:46 +0100 CET", date, "")
}
