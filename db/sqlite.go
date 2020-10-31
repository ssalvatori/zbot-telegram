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

	log "github.com/sirupsen/logrus"
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
	log.Debug("Closing connection")
	d.Close()
}

//Init connect to sqlite db
func (d *ZbotDatabaseSqlite) Init() error {
	log.Info("Connecting to database: " + d.File)

	newLogger := logger.New(
		log.New(), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(d.File), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now()
		},
	})

	if err != nil {
		log.Error(err)
		return err
	}
	if db == nil {
		log.Error(err)
		return errors.New("Error connecting")
	}

	db.Debug().AutoMigrate(&Definition{}, &UserIgnore{})

	d.DB = db

	return nil
}

//Statistics get total number of definitions
func (d *ZbotDatabaseSqlite) Statistics() (string, error) {

	var count int64
	if result := d.DB.Model(&Definition{}).Count(&count); result.Error != nil {
		log.Error(result.Error)
		return "", result.Error
	}

	return strconv.FormatInt(count, 10), nil
}

//Last get last X definitions
func (d *ZbotDatabaseSqlite) Last(limit int) ([]Definition, error) {
	definitions := []Definition{}
	if err := d.DB.Debug().Model(&Definition{}).Limit(limit).Order("ID desc").Find(&definitions).Error; err != nil {
		log.Error(err)
		return nil, err
	}
	return definitions, nil
}

//Top get definition with the most numbers of hits
func (d *ZbotDatabaseSqlite) Top(limit int) ([]Definition, error) {
	definitions := []Definition{}
	if result := d.DB.Debug().Model(&Definition{}).Limit(limit).Order("hits desc").Find(&definitions); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Definition{}, ErrNotFound
		}
		return []Definition{}, result.Error
	}
	return definitions, nil
}

//Rand get a random definition from the DB
func (d *ZbotDatabaseSqlite) Rand(limit int) ([]Definition, error) {
	definitions := []Definition{}
	if result := d.DB.Debug().Model(&Definition{}).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
		log.Error(result.Error)
		return []Definition{}, result.Error
	}

	return definitions, nil
}

//Get get a definition from the db using term
func (d *ZbotDatabaseSqlite) Get(term string, chat string) (Definition, error) {
	var def Definition

	if result := d.DB.Debug().Model(&Definition{}).Where("term = ? COLLATE NOCASE", term).First(&def); result.Error != nil {
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
		log.Error(err)
		return ErrInternalError
	}
	return nil
}

//Find terms with criteria inside of the meaning
func (d *ZbotDatabaseSqlite) Find(criteria string, chat string, limit int) ([]Definition, error) {

	definitions := []Definition{}
	criteria = fmt.Sprintf("%%%s%%", criteria)

	if result := d.DB.Debug().Model(&Definition{}).Where("meaning like ? COLLATE NOCASE", criteria).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
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

	if result := d.DB.Debug().Model(&Definition{}).Where("term like ? COLLATE NOCASE", criteria).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
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
	log.Debug(definition)
	for {
		err := d._set(term, definition)
		if err != nil {
			log.Debug("SQL insert error: ", err.Error())
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				term = fmt.Sprintf("%s%d", definition.Term, count)
				log.Debug(fmt.Sprintf("New Term: %s", term))
				count = count + 1
			} else {
				return "", err
			}
		} else {
			log.Debug("trying with: ", term)
			break
		}
	}
	return term, nil
}

//_set create a new definition
func (d *ZbotDatabaseSqlite) _set(term string, definition Definition) error {

	definition.Term = term

	if err := d.DB.Debug().Model(&Definition{}).Create(&definition).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//Append append meaning to a given definition
func (d *ZbotDatabaseSqlite) Append(item Definition) error {

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
			log.Error(err)
			return err
		}
		return nil
	}
	return ErrLocked
}

//Lock a given definition (no more append or forget)
func (d *ZbotDatabaseSqlite) Lock(item Definition) error {
	definition, err := d.Get(item.Term, item.Chat)
	if err != nil {
		log.Error(err)
		return err
	}

	if !definition.Locked.Bool {
		def := Definition{Locked: sql.NullBool{Bool: true, Valid: true}, LockedBy: sql.NullString{String: item.Author}}
		if err := d.DB.Debug().Model(&definition).Updates(def).Error; err != nil {
			log.Error(err)
			return err
		}
		return nil
	}

	return fmt.Errorf("Already locked by  %q", definition.LockedBy.String)
}

func (d *ZbotDatabaseSqlite) Forget(item Definition) error {
	return nil
}

func (d *ZbotDatabaseSqlite) UserIgnoreList() ([]UserIgnore, error) {
	return nil, nil
}

//UserLevel Mock
func (d *ZbotDatabaseSqlite) UserLevel(str string) (string, error) {
	return "bnil", nil
}

//UserCheckIgnore Mock, it will return false if error is set otherwise it will return Ignore_User value
func (d *ZbotDatabaseSqlite) UserCheckIgnore(str string) bool {
	return true
}

func (d *ZbotDatabaseSqlite) UserIgnoreInsert(username string) error {
	return nil
}

//UserCleanupIgnorelist Cleanup ignore list
func (d *ZbotDatabaseSqlite) UserCleanupIgnorelist() error {
	return nil
}
