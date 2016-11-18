package main

import (
	"database/sql"
)

type mockZbotDatabase struct {
	level string
	file string
}

func (d mockZbotDatabase) init() error {
	return nil
}

func (d mockZbotDatabase) close() {
}

func (d mockZbotDatabase) statistics() (string, error) {
	return d.level, nil
}

func (d mockZbotDatabase) top() ([]definitionItem, error) {
	var items []definitionItem
	return items, nil
}

func (d mockZbotDatabase) rand() (definitionItem, error) {
	var item definitionItem
	return item, nil
}

func (d mockZbotDatabase) last() (definitionItem, error) {
	return definitionItem{},nil
}

func (d mockZbotDatabase) set(def definitionItem) error {
	return nil
}

func (d mockZbotDatabase) find(criteria string) ([]definitionItem, error) {
	var items []definitionItem
	return items, nil
}

func (d mockZbotDatabase) get(term string) (definitionItem, error) {
	var item definitionItem
	return item, nil
}

func (d mockZbotDatabase) _set(term string, def definitionItem) (sql.Result, error) {
	var result sql.Result
	return result, nil
}

func (d mockZbotDatabase) search(str string) ([]definitionItem, error) {
	var def []definitionItem
	return def, nil
}

func (d mockZbotDatabase) userLevel(str string) (string, error) {
	var strr string
	return strr, nil
}

func (d mockZbotDatabase) userCheckIgnore(str string) (bool, error) {
	var ignore = false
	return ignore, nil
}

func (d mockZbotDatabase) userIgnoreInsert(username string) error {
	return nil
}

func (d mockZbotDatabase) userCleanIgnore() error {
	return nil
}
