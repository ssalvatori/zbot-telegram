package database

import "database/sql"

type definitionItem struct {
	term    string
	meaning string
	author  string
	date    string
	id      int
}

type userIgnore struct {
	username string
	since string
	until string
}

type zbotDatabase interface {
	init() error
	close()
	statistics() (string, error)
	top() ([]definitionItem, error)
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
	userIgnoreList() ([]userIgnore, error)
}



