package utils

import (
	"errors"
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

	assert.Equal(t, "2020-11-01 10:10:46 +0000 UTC", date, "")
}

func TestStringToArray(t *testing.T) {
	assert.Equal(t, []string{"cmd1", "cmd2", "cmd3", "cmd4"}, StringToArray("cmd1,cmd2,cmd3, cmd4"), "Set Disabled Commands")
	assert.Equal(t, []string{}, StringToArray(""), "No commands")
}

func TestParseCommand(t *testing.T) {

	tests := []struct {
		in          string
		out         string
		outErr      error
		description string
	}{
		{"/cmd", "cmd", nil, "cmd without arg"},
		{"/cmd arg1", "cmd", nil, "cmd with single arg"},
		{"/cmd arg1 arg2", "cmd", nil, "cmd with multiples args"},
		{"arg1 arg2 arg3", "", errors.New("Text could not been parser"), "without cmd"},
	}

	for _, test := range tests {
		got, err := ParseCommand(test.in)
		assert.Equal(t, test.out, got, test.description)
		assert.Equal(t, test.outErr, err, test.description)
	}
}

func TestGetCommandFile(t *testing.T) {
	var cmdList = []struct {
		Key         string
		File        string
		Description string
	}{
		{File: "cmdFile1", Key: "cmd1", Description: ""},
		{File: "cmdFile2", Key: "cmd2", Description: ""},
	}

	file, err := GetCommandFile("cmd1", cmdList)
	assert.Equal(t, "cmdFile1", file, "key exits")
	assert.Equal(t, nil, err, "key exits")

	file, err = GetCommandFile("cmd2", cmdList)
	assert.Equal(t, "cmdFile2", file, "key exits")
	assert.Equal(t, nil, err, "key exits")

	file, err = GetCommandFile("cmd3", cmdList)
	assert.Equal(t, "", file, "key doesn't exits")
	assert.Error(t, err, "Key doesn't exits")
}
