package db

import "database/sql"

//ZbotDatabase DB interface for Zbot
type ZbotDatabase interface {
	GetConnectionInfo() string
	Init() error
	Close()
	Statistics() (string, error)
	Append(DefinitionItem) error
	Top() ([]DefinitionItem, error)
	Rand() (DefinitionItem, error)
	Last() (DefinitionItem, error)
	Get(string) (DefinitionItem, error)
	Set(DefinitionItem) (string, error)
	_set(string, DefinitionItem) (sql.Result, error)
	Find(string) ([]DefinitionItem, error)
	Search(string) ([]DefinitionItem, error)
	Forget(DefinitionItem) error
	UserLevel(string) (string, error)
	UserIgnoreInsert(string) error
	//UserCheckIgnore return true if the user is on the ignore_list, false if it isn´t
	UserCheckIgnore(string) bool
	UserCleanIgnore() error
	UserIgnoreList() ([]UserIgnore, error)

	Lock(DefinitionItem) error
}

//DefinitionItem .
type DefinitionItem struct {
	Term    string
	Meaning string
	Author  string
	Date    string
	Id      int
}

//UserIgnore .
type UserIgnore struct {
	Username string
	Since    string
	Until    string
}
