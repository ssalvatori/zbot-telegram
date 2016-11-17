package main

import (
	"database/sql"
)

type mockDatabase struct {
}

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

func (d mockDatabase) rand() (definitionItem, error) {
	var item definitionItem
	return item, nil
}

func (d mockDatabase) get(term string) (definitionItem, error) {
	var item definitionItem
	return item, nil
}

func (d mockDatabase) _set(term string, def definitionItem) (sql.Result, error) {
	var result sql.Result
	return result, nil
}

func (d mockDatabase) search(str string) ([]definitionItem, error) {
	var def []definitionItem
	return def, nil
}

func (d mockDatabase) userLevel(str string) (string, error) {
	var strr string
	return strr, nil
}

func (d mockDatabase) userCheckIgnore(str string) (bool, error) {
	var ignore = false
	return ignore, nil
}

func (d mockDatabase) userCleanIgnore() error {
	return nil
}
