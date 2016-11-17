package main

import (
	"testing"
)

type mockDatabase struct {
}

/*
rand() (definitionItem, error)
last() (definitionItem, error)
get(string) (definitionItem, error)
set(definitionItem) error
_set(string, definitionItem) (sql.Result, error)
find(string) ([]definitionItem, error)
search(string) ([]definitionItem, error)
userLevel(string) (string, error)
userIgnoreInsert(string) error
userCheckIgnore(string) (bool, error)
userCleanIgnore() error
*/

func (d mockDatabase) init() error {
	return nil
}

func (d mockDatabase) close() {
}

func (d mockDatabase) statistics() (string, error) {
	return "66", nil
}

func (d mockDatabase) top() ([]definitionItem, error) {
	var items []definitionItem
	return items, nil
}