package main

import (
	"database/sql"
	"github.com/tucnak/telebot"
)

type mockZbotDatabase struct {
	level string
	file string
	term string
	meaning string
	find_terms []string
	search_terms []string
	rand_def definitionItem
	user telebot.User
	ignore_list []string
	user_ignored []userIgnore
}

func (d *mockZbotDatabase) init() error {
	return nil
}

func (d *mockZbotDatabase) close() {
}

func (d *mockZbotDatabase) userIgnoreList() ([]userIgnore, error) {
	return d.user_ignored, nil
}

func (d *mockZbotDatabase) statistics() (string, error) {
	return d.level, nil
}

func (d *mockZbotDatabase) top() ([]definitionItem, error) {
	var items []definitionItem

	for _,value := range d.find_terms {
		items = append(items, definitionItem{term: value})
	}

	return items, nil
}

func (d *mockZbotDatabase) rand() (definitionItem, error) {
	return d.rand_def, nil
}

func (d *mockZbotDatabase) last() (definitionItem, error) {
	return definitionItem{term: d.term, meaning:d.meaning},nil
}

func (d *mockZbotDatabase) set(def definitionItem) error {
	return nil
}

func (d *mockZbotDatabase) find(criteria string) ([]definitionItem, error) {
	return []definitionItem{definitionItem{term: d.term}}, nil
}

func (d *mockZbotDatabase) get(term string) (definitionItem, error) {
	return definitionItem{term: d.term, meaning:d.meaning}, nil
}

func (d *mockZbotDatabase) _set(term string, def definitionItem) (sql.Result, error) {
	var result sql.Result
	return result, nil
}

func (d *mockZbotDatabase) search(str string) ([]definitionItem, error) {
	var items []definitionItem

	for _,value := range d.search_terms {
		items = append(items, definitionItem{term: value})
	}


	return items, nil
}

func (d *mockZbotDatabase) userLevel(str string) (string, error) {
	return d.level, nil
}

func (d *mockZbotDatabase) userCheckIgnore(str string) (bool, error) {
	var ignore = false
	return ignore, nil
}

func (d *mockZbotDatabase) userIgnoreInsert(username string) error {
	return nil
}

func (d *mockZbotDatabase) userCleanIgnore() error {
	return nil
}
