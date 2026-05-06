package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log/slog"
)

//ZbotDatabaseSqlite struct
type ZbotDatabaseSqlite struct {
	DB   *gorm.DB
	File string
}

//GetConnectionInfo get connection information
func (d *ZbotDatabaseSqlite) GetConnectionInfo() string {
	return fmt.Sprintf("DB: %s", d.File)
}

//Close close connecttion to DB
func (d *ZbotDatabaseSqlite) Close() {
	slog.Debug("Closing DB connection")
}

//Init connect to sqlite db
func (d *ZbotDatabaseSqlite) Init() error {
	slog.Info("Connecting to database", "file", d.File)

	newLogger := logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(sqlite.Open(d.File), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now()
		},
	})

	if err != nil {
		slog.Error("open db error", "err", err)
		return err
	}
	if db == nil {
		slog.Error("nil db after open")
		return errors.New("Error connecting")
	}

	err = db.Debug().AutoMigrate(&Definition{}, &UserIgnore{})
	if err != nil {
		slog.Error("migration error", "err", err)
		return errors.New("Error during migration")
	}

	d.DB = db

	return nil
}

//Statistics get total number of definitions
func (d *ZbotDatabaseSqlite) Statistics(chat string) (string, error) {

	var count int64
	if result := d.DB.Model(&Definition{}).Where("chat = ? COLLATE NOCASE", chat).Count(&count); result.Error != nil {
		slog.Error("statistics query error", "err", result.Error)
		return "", result.Error
	}

	return strconv.FormatInt(count, 10), nil
}

//Last get last X definitions
func (d *ZbotDatabaseSqlite) Last(chat string, limit int) ([]Definition, error) {
	definitions := []Definition{}
	if err := d.DB.Debug().Model(&Definition{}).Where("chat = ? COLLATE NOCASE", chat).Limit(limit).Order("ID desc").Find(&definitions).Error; err != nil {
		slog.Error("last query error", "err", err)
		return nil, err
	}
	return definitions, nil
}

//Top get definition with the most numbers of hits
func (d *ZbotDatabaseSqlite) Top(chat string, limit int) ([]Definition, error) {
	definitions := []Definition{}
	if result := d.DB.Debug().Model(&Definition{}).Where("chat = ? COLLATE NOCASE", chat).Limit(limit).Order("hits desc").Find(&definitions); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Definition{}, ErrNotFound
		}
		return []Definition{}, result.Error
	}
	return definitions, nil
}

//Rand get a random definition from the DB
func (d *ZbotDatabaseSqlite) Rand(chat string, limit int) ([]Definition, error) {
	definitions := []Definition{}
	if result := d.DB.Debug().Model(&Definition{}).Where("chat = ? COLLATE NOCASE", chat).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
		slog.Error("rand query error", "err", result.Error)
		return []Definition{}, result.Error
	}

	return definitions, nil
}

//Get get a definition from the db using term
func (d *ZbotDatabaseSqlite) Get(term string, chat string) (Definition, error) {
	var def Definition

	if result := d.DB.Debug().Model(&Definition{}).Where("term = ? COLLATE NOCASE", term, chat).Where("chat = ? COLLATE NOCASE", chat).First(&def); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Definition{}, ErrNotFound
		}
		return Definition{}, result.Error
	}

	return def, nil
}

//IncreaseHits increase the definition hits by one
func (d *ZbotDatabaseSqlite) IncreaseHits(id uint) error {
	if err := d.DB.Debug().Model(&Definition{}).Where("id = ?", id).UpdateColumn("hits", gorm.Expr("hits + ?", 1)).Error; err != nil {
		slog.Error("increase hits error", "err", err)
		return ErrInternalError
	}
	return nil
}

//Find terms with criteria inside of the meaning
func (d *ZbotDatabaseSqlite) Find(criteria string, chat string, limit int) ([]Definition, error) {

	definitions := []Definition{}
	criteria = fmt.Sprintf("%%%s%%", criteria)

	if result := d.DB.Debug().Model(&Definition{}).Where("meaning like ? COLLATE NOCASE", criteria).Where("chat = ? COLLATE NOCASE", chat).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Definition{}, ErrNotFound
		}
	} else {
		if result.RowsAffected == 0 {
			return []Definition{}, ErrNotFound
		}
	}
	return definitions, nil
}

//Search find list of term with a given pattern
func (d *ZbotDatabaseSqlite) Search(criteria string, chat string, limit int) ([]Definition, error) {
	definitions := []Definition{}
	criteria = fmt.Sprintf("%%%s%%", criteria)

	if result := d.DB.Debug().Model(&Definition{}).Where("term like ? COLLATE NOCASE", criteria).Where("chat = ? COLLATE NOCASE", chat).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Definition{}, ErrNotFound
		}
	} else {
		if result.RowsAffected == 0 {
			return []Definition{}, ErrNotFound
		}
	}
	return definitions, nil
}

//Set create a new term in the db.
func (d *ZbotDatabaseSqlite) Set(definition Definition) (string, error) {
	count := 1
	term := definition.Term
	slog.Debug("set definition", "definition", definition)
	for {
		err := d._set(term, definition)
		if err != nil {
			slog.Debug("SQL insert error", "err", err.Error())
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				term = fmt.Sprintf("%s%d", definition.Term, count)
				slog.Debug("new term", "term", term)
				count = count + 1
			} else {
				return "", err
			}
		} else {
			slog.Debug("trying with term", "term", term)
			break
		}
	}
	return term, nil
}

//_set create a new definition
func (d *ZbotDatabaseSqlite) _set(term string, definition Definition) error {

	definition.Term = term

	if err := d.DB.Debug().Model(&Definition{}).Create(&definition).Error; err != nil {
		slog.Error("create definition error", "err", err)
		return err
	}
	return nil
}

//Append append meaning to a given definition
func (d *ZbotDatabaseSqlite) Append(item Definition, chat string) error {

	definition, err := d.Get(item.Term, item.Chat)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fmt.Errorf("Definition [%s] not found", item.Term)
		}
		return ErrInternalError
	}

	if !definition.Locked.Bool {
		appenedMeaning := fmt.Sprintf("%s. %s", definition.Meaning, item.Meaning)
		def := Definition{Meaning: appenedMeaning, Author: item.Author}
		if err := d.DB.Debug().Model(&definition).Updates(def).Error; err != nil {
			slog.Error("append definition error", "err", err)
			return err
		}
		return nil
	}
	return ErrLocked
}

//Lock a given definition (no more append or forget)
func (d *ZbotDatabaseSqlite) Lock(item Definition, chat string) error {
	definition, err := d.Get(item.Term, item.Chat)
	if err != nil {
		slog.Error("lock get error", "err", err)
		return err
	}

	if !definition.Locked.Bool {
		def := Definition{Locked: sql.NullBool{Bool: true, Valid: true}, LockedBy: sql.NullString{String: item.Author}}
		if err := d.DB.Debug().Model(&definition).Updates(def).Where("chat = ? COLLATE NOCASE", chat).Error; err != nil {
			slog.Error("lock update error", "err", err)
			return err
		}
		return nil
	}

	return fmt.Errorf("Already locked by  %q", definition.LockedBy.String)
}

//Forget No implemented yet
func (d *ZbotDatabaseSqlite) Forget(item Definition, chat string) error {
	return nil
}

//UserIgnoreList No Implemeted yet
func (d *ZbotDatabaseSqlite) UserIgnoreList() ([]UserIgnore, error) {
	return nil, nil
}

//UserLevel Mock
func (d *ZbotDatabaseSqlite) UserLevel(str string) (string, error) {
	return "bnil", nil
}

//UserCheckIgnore Mock, it will return false if error is set otherwise it will return IgnoreUser value
func (d *ZbotDatabaseSqlite) UserCheckIgnore(str string) bool {
	return true
}

//UserIgnoreInsert Add a user to ignore list
func (d *ZbotDatabaseSqlite) UserIgnoreInsert(username string) error {
	return nil
}

//UserCleanupIgnorelist Cleanup ignore list
func (d *ZbotDatabaseSqlite) UserCleanupIgnorelist() error {
	return nil
}
