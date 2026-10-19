package db

import (
	"errors"

	"gopkg.in/telebot.v4"
)

// ZbotDatabaseMock mock object
type ZbotDatabaseMock struct {
	Level             string
	File              string
	Term              string
	Meaning           string
	Author            string
	CreateAt          int64
	UpdateAt          int64
	FindTerms         []string
	SearchTerms       []string
	NotFound          bool
	RandDef           []Definition
	User              telebot.User
	IgnoreList        []string
	UserIgnored       []UserIgnore
	IgnoreUser        bool
	Error             bool
	IgnoreListCleaned bool
	Channel           Channel
}

// Engine mock
func (d *ZbotDatabaseMock) Engine() string {
	return "mock"
}

// GetConnectionInfo mock
func (d *ZbotDatabaseMock) GetConnectionInfo() string {
	return "mock"
}

// Init mock
func (d *ZbotDatabaseMock) Init() error {
	return nil
}

// Close mock
func (d *ZbotDatabaseMock) Close() {
}

// Statistics mock
func (d *ZbotDatabaseMock) Statistics(channelID uint) (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return d.Level, nil
}

// Top mock
func (d *ZbotDatabaseMock) Top(channelID uint, limit int) ([]Definition, error) {
	var items []Definition

	if d.Error {
		return nil, errors.New("mock")
	}

	for _, value := range d.FindTerms {
		items = append(items, Definition{Term: value})
	}

	return items, nil
}

// Rand mock
func (d *ZbotDatabaseMock) Rand(channelID uint, limit int) ([]Definition, error) {
	if d.Error {
		return []Definition{}, errors.New("mock")
	}
	return d.RandDef, nil
}

// Last mock
func (d *ZbotDatabaseMock) Last(channelID uint, last int) ([]Definition, error) {

	if d.Error {
		return []Definition{}, errors.New("mock")
	}

	return []Definition{{Term: d.Term, Meaning: d.Meaning}}, nil
}

// Set mock
func (d *ZbotDatabaseMock) Set(def Definition) (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return def.Term, nil
}

// Find mock
func (d *ZbotDatabaseMock) Find(criteria string, channelID uint, limit int) ([]Definition, error) {
	if d.NotFound {
		return []Definition{}, nil
	}
	if d.Error {
		return nil, errors.New("mock")
	}
	return []Definition{{Term: d.Term}}, nil
}

// Get definition mock
func (d *ZbotDatabaseMock) Get(term string, channelID uint) (Definition, error) {
	if d.NotFound {
		return Definition{}, ErrNotFound
	}
	if d.Error {
		return Definition{}, errors.New("mock")
	}
	return Definition{Term: d.Term, Meaning: d.Meaning, Author: d.Author, UpdatedAt: d.UpdateAt}, nil
}

// _set mock
func (d *ZbotDatabaseMock) _set(term string, def Definition) error {

	if d.Error {
		return errors.New("mock")
	}

	return nil
}

// Search mock
func (d *ZbotDatabaseMock) Search(str string, channelID uint, limit int) ([]Definition, error) {
	var items []Definition

	if d.Error {
		return []Definition{}, errors.New("mock")
	}

	for _, value := range d.SearchTerms {
		items = append(items, Definition{Term: value})
	}

	return items, nil
}

// UserLevel Mock
func (d *ZbotDatabaseMock) UserLevel(str string) (string, error) {
	if d.Error {
		return "", errors.New("mock")
	}
	return d.Level, nil
}

// UserCheckIgnore Mock, it will return false if error is set otherwise it will return IgnoreUser value
func (d *ZbotDatabaseMock) UserCheckIgnore(str string, channelID uint) bool {

	if d.Error {
		return false
	}

	return d.IgnoreUser
}

// UserIgnoreInsert mock
func (d *ZbotDatabaseMock) UserIgnoreInsert(username string, durationSeconds int, createdBy string, channelID uint) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}

// UserCleanupIgnorelist mock
func (d *ZbotDatabaseMock) UserCleanupIgnorelist() error {
	d.IgnoreListCleaned = true
	return nil
}

// Lock mock
func (d *ZbotDatabaseMock) Lock(item Definition) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}

// Append mock
func (d *ZbotDatabaseMock) Append(item Definition) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}

// Forget mock
func (d *ZbotDatabaseMock) Forget(item Definition) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}

// UserIgnoreList mock
func (d *ZbotDatabaseMock) UserIgnoreList() ([]UserIgnore, error) {
	if d.Error {
		return nil, errors.New("mock")
	}
	return d.UserIgnored, nil
}

// IncreaseHits mock
func (d *ZbotDatabaseMock) IncreaseHits(limit uint) error {
	return nil
}

func (d *ZbotDatabaseMock) AddChannel(telegramID int64, name string) error {
	if d.Error {
		return errors.New("mock")
	}
	return nil
}

func (d *ZbotDatabaseMock) RemoveChannel(telegramID int64) (bool, error) {
	if d.Error {
		return false, errors.New("mock")
	}
	return true, nil
}

func (d *ZbotDatabaseMock) ListChannels() ([]Channel, error) {
	if d.Error {
		return nil, errors.New("mock")
	}
	return []Channel{}, nil
}

// GetChannel returns d.Channel, defaulting to a channel with ID 1 when unset.
func (d *ZbotDatabaseMock) GetChannel(telegramID int64) (Channel, error) {
	if d.Error {
		return Channel{}, errors.New("mock")
	}
	if d.Channel.ID == 0 {
		return Channel{ID: 1, TelegramID: telegramID, Name: "mock", LearnEnabled: true, ExternalEnabled: true}, nil
	}
	return d.Channel, nil
}
