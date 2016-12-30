package db

import (
	"database/sql"
	"github.com/tucnak/telebot"
)

type MockZbotDatabase struct {
	Level        string
	File         string
	Term         string
	Meaning      string
	Find_terms   []string
	Search_terms []string
	Rand_def     DefinitionItem
	User         telebot.User
	Ignore_list  []string
	User_ignored []UserIgnore
}

func (d *MockZbotDatabase) Init() error {
	return nil
}

func (d *MockZbotDatabase) Close() {
}

func (d *MockZbotDatabase) UserIgnoreList() ([]UserIgnore, error) {
	return d.User_ignored, nil
}

func (d *MockZbotDatabase) Statistics() (string, error) {
	return d.Level, nil
}

func (d *MockZbotDatabase) Top() ([]DefinitionItem, error) {
	var items []DefinitionItem

	for _,value := range d.Find_terms {
		items = append(items, DefinitionItem{Term: value})
	}

	return items, nil
}

func (d *MockZbotDatabase) Rand() (DefinitionItem, error) {
	return d.Rand_def, nil
}

func (d *MockZbotDatabase) Last() (DefinitionItem, error) {
	return DefinitionItem{Term: d.Term, Meaning:d.Meaning},nil
}

func (d *MockZbotDatabase) Set(def DefinitionItem) error {
	return nil
}

func (d *MockZbotDatabase) Find(criteria string) ([]DefinitionItem, error) {
	return []DefinitionItem{DefinitionItem{Term: d.Term}}, nil
}

func (d *MockZbotDatabase) Get(term string) (DefinitionItem, error) {
	return DefinitionItem{Term: d.Term, Meaning:d.Meaning}, nil
}

func (d *MockZbotDatabase) _set(term string, def DefinitionItem) (sql.Result, error) {
	var result sql.Result
	return result, nil
}

func (d *MockZbotDatabase) Search(str string) ([]DefinitionItem, error) {
	var items []DefinitionItem

	for _,value := range d.Search_terms {
		items = append(items, DefinitionItem{Term: value})
	}


	return items, nil
}

func (d *MockZbotDatabase) UserLevel(str string) (string, error) {
	return d.Level, nil
}

func (d *MockZbotDatabase) UserCheckIgnore(str string) (bool, error) {
	var ignore = false
	return ignore, nil
}

func (d *MockZbotDatabase) UserIgnoreInsert(username string) error {
	return nil
}

func (d *MockZbotDatabase) UserCleanIgnore() error {
	return nil
}
