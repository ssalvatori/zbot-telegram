package db

import (
	"database/sql"

	"errors"

	"gopkg.in/tucnak/telebot.v2"
)

//ZbotDatabaseMock mock object
type ZbotDatabaseMock struct {
	Level             string
	File              string
	Term              string
	Meaning           string
	Author            string
	Date              string
	FindTerms         []string
	Search_terms      []string
	Not_found         bool
	Rand_def          Definition
	User              telebot.User
	IgnoreList        []string
	UserIgnored       []UserIgnore
	IgnoreUser        bool
	Error             bool
	ErrorAppend       bool
	IgnoreListCleaned bool
}

func (d *ZbotDatabaseMock) GetConnectionInfo() string {
	return "mock"
}

func (d *ZbotDatabaseMock) Init() error {
	return nil
}

func (d *ZbotDatabaseMock) Close() {
}

func (d *ZbotDatabaseMock) UserIgnoreList() ([]UserIgnore, error) {
	if d.Error {
		return nil, errors.New("mock")
	}
	return d.UserIgnored, nil
}

//Statistics mock
func (d *ZbotDatabaseMock) Statistics(string) (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return d.Level, nil
}

func (d *ZbotDatabaseMock) Top() ([]Definition, error) {
	var items []Definition

	if d.Error {
		return nil, errors.New("mock")
	}

	for _, value := range d.FindTerms {
		items = append(items, Definition{Term: value})
	}

	return items, nil
}

func (d *ZbotDatabaseMock) Rand() (Definition, error) {
	if d.Error {
		return Definition{}, errors.New("mock")
	}
	return d.Rand_def, nil
}

func (d *ZbotDatabaseMock) Last() (Definition, error) {

	if d.Error {
		return Definition{}, errors.New("mock")
	}

	return Definition{Term: d.Term, Meaning: d.Meaning}, nil
}

func (d *ZbotDatabaseMock) Set(def Definition) (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return def.Term, nil
}

//Find mock to find terms looking inside of the meaning text
func (d *ZbotDatabaseMock) Find(criteria string, chat string) ([]Definition, error) {
	if d.Not_found {
		return []Definition{}, nil
	}
	if d.Error {
		return nil, errors.New("mock")
	}
	return []Definition{{Term: d.Term}}, nil
}

//Get definition mock
func (d *ZbotDatabaseMock) Get(term string, chat string) (Definition, error) {
	if d.Not_found {
		return Definition{}, nil
	}
	if d.Error {
		return Definition{}, errors.New("mock")
	}
	return Definition{Term: d.Term, Meaning: d.Meaning, Author: d.Author, Date: d.Date}, nil
}

func (d *ZbotDatabaseMock) _set(term string, def Definition) (sql.Result, error) {
	var result sql.Result

	if d.Error {
		return nil, errors.New("mock")
	}

	return result, nil
}

//Search mock
func (d *ZbotDatabaseMock) Search(str string, chat string) ([]Definition, error) {
	var items []Definition

	if d.Error {
		return []Definition{}, errors.New("mock")
	}

	for _, value := range d.Search_terms {
		items = append(items, Definition{Term: value})
	}

	return items, nil
}

//UserLevel Mock
func (d *ZbotDatabaseMock) UserLevel(str string) (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return d.Level, nil
}

//UserCheckIgnore Mock, it will return false if error is set otherwise it will return IgnoreUser value
func (d *ZbotDatabaseMock) UserCheckIgnore(str string) bool {

	if d.Error {
		return false
	}

	return d.IgnoreUser
}

func (d *ZbotDatabaseMock) UserIgnoreInsert(username string) error {
	return nil
}

//UserCleanupIgnorelist Cleanup ignore list
func (d *ZbotDatabaseMock) UserCleanupIgnorelist() error {
	d.IgnoreListCleaned = true
	return nil
}

//Lock defintion mock
func (d *ZbotDatabaseMock) Lock(item Definition) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}

func (d *ZbotDatabaseMock) Append(item Definition) error {
	if d.ErrorAppend {
		return errors.New("mock")
	}
	return nil
}

func (d *ZbotDatabaseMock) Forget(item Definition) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}
