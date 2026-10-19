package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log/slog"
)

// ZbotDatabaseSqlite struct
type ZbotDatabaseSqlite struct {
	DB   *gorm.DB
	File string
}

func newGORMLogger(applicationLogger *slog.Logger) logger.Interface {
	logLevel := logger.Silent
	if applicationLogger.Enabled(context.Background(), slog.LevelDebug) {
		logLevel = logger.Info
	}
	return logger.NewSlogLogger(applicationLogger, logger.Config{LogLevel: logLevel})
}

// Engine returns the backend name
func (d *ZbotDatabaseSqlite) Engine() string {
	return EngineSqlite
}

// GetConnectionInfo get connection information
func (d *ZbotDatabaseSqlite) GetConnectionInfo() string {
	return fmt.Sprintf("DB: %s", d.File)
}

// Close close connecttion to DB
func (d *ZbotDatabaseSqlite) Close() {
	slog.Debug("Closing DB connection")
	if d.DB == nil {
		return
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		slog.Error("get SQL database connection error", "err", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		slog.Error("close database connection error", "err", err)
	}
}

// Init connect to sqlite db
func (d *ZbotDatabaseSqlite) Init() error {
	slog.Info("Connecting to database", "file", d.File)

	db, err := gorm.Open(sqlite.Open(d.File), &gorm.Config{
		Logger:         newGORMLogger(slog.Default()),
		TranslateError: true,
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
		return errors.New("error connecting")
	}

	err = db.AutoMigrate(&Definition{}, &UserIgnore{}, &Channel{})
	if err != nil {
		slog.Error("migration error", "err", err)
		return errors.New("error during migration")
	}

	d.DB = db

	return nil
}

// Statistics get total number of definitions
func (d *ZbotDatabaseSqlite) Statistics(channelID uint) (string, error) {

	var count int64
	if result := d.DB.Model(&Definition{}).Where("channel_id = ?", channelID).Count(&count); result.Error != nil {
		slog.Error("statistics query error", "err", result.Error)
		return "", result.Error
	}

	return strconv.FormatInt(count, 10), nil
}

// Last get last X definitions
func (d *ZbotDatabaseSqlite) Last(channelID uint, limit int) ([]Definition, error) {
	definitions := []Definition{}
	if err := d.DB.Model(&Definition{}).Where("channel_id = ?", channelID).Limit(limit).Order("ID desc").Find(&definitions).Error; err != nil {
		slog.Error("last query error", "err", err)
		return nil, err
	}
	return definitions, nil
}

// Top get definition with the most numbers of hits
func (d *ZbotDatabaseSqlite) Top(channelID uint, limit int) ([]Definition, error) {
	definitions := []Definition{}
	if result := d.DB.Model(&Definition{}).Where("channel_id = ?", channelID).Limit(limit).Order("hits desc").Find(&definitions); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Definition{}, ErrNotFound
		}
		return []Definition{}, result.Error
	}
	return definitions, nil
}

// Rand get a random definition from the DB
func (d *ZbotDatabaseSqlite) Rand(channelID uint, limit int) ([]Definition, error) {
	definitions := []Definition{}
	if result := d.DB.Model(&Definition{}).Where("channel_id = ?", channelID).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
		slog.Error("rand query error", "err", result.Error)
		return []Definition{}, result.Error
	}

	return definitions, nil
}

// Get get a definition from the db using term
func (d *ZbotDatabaseSqlite) Get(term string, channelID uint) (Definition, error) {
	var def Definition

	if result := d.DB.Model(&Definition{}).Where("term = ? COLLATE NOCASE", term).Where("channel_id = ?", channelID).First(&def); result.Error != nil {
		slog.Error("get query error", "err", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Definition{}, ErrNotFound
		}
		return Definition{}, result.Error
	}

	return def, nil
}

// IncreaseHits increase the definition hits by one
func (d *ZbotDatabaseSqlite) IncreaseHits(id uint) error {
	if err := d.DB.Model(&Definition{}).Where("id = ?", id).UpdateColumn("hits", gorm.Expr("hits + ?", 1)).Error; err != nil {
		slog.Error("increase hits error", "err", err)
		return ErrInternalError
	}
	return nil
}

// Find terms with criteria inside of the meaning
func (d *ZbotDatabaseSqlite) Find(criteria string, channelID uint, limit int) ([]Definition, error) {

	definitions := []Definition{}
	criteria = fmt.Sprintf("%%%s%%", criteria)

	if result := d.DB.Model(&Definition{}).Where("meaning like ? COLLATE NOCASE", criteria).Where("channel_id = ?", channelID).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
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

// Search find list of term with a given pattern
func (d *ZbotDatabaseSqlite) Search(criteria string, channelID uint, limit int) ([]Definition, error) {
	definitions := []Definition{}
	criteria = fmt.Sprintf("%%%s%%", criteria)

	if result := d.DB.Model(&Definition{}).Where("term like ? COLLATE NOCASE", criteria).Where("channel_id = ?", channelID).Limit(limit).Order("random()").Find(&definitions); result.Error != nil {
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

const maxSetAttempts = 100

// Set create a new term in the db.
func (d *ZbotDatabaseSqlite) Set(definition Definition) (string, error) {
	slog.Debug("set definition", "definition", definition)
	for attempt := 0; attempt < maxSetAttempts; attempt++ {
		term := definition.Term
		if attempt > 0 {
			term = fmt.Sprintf("%s%d", definition.Term, attempt)
		}

		err := d._set(term, definition)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				slog.Debug("term already exists, trying next suffix", "term", term)
				continue
			}
			slog.Error("create definition error", "err", err)
			return "", err
		}

		return term, nil
	}

	return "", fmt.Errorf("%w: %q after %d attempts", ErrTermCollisionLimit, definition.Term, maxSetAttempts)
}

// _set create a new definition
func (d *ZbotDatabaseSqlite) _set(term string, definition Definition) error {

	definition.Term = term

	if err := d.DB.Model(&Definition{}).Create(&definition).Error; err != nil {
		return err
	}
	return nil
}

// Append append meaning to a given definition
func (d *ZbotDatabaseSqlite) Append(item Definition) error {

	definition, err := d.Get(item.Term, item.ChannelID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fmt.Errorf("Definition [%s] not found", item.Term)
		}
		return ErrInternalError
	}

	if !definition.Locked.Bool {
		appenedMeaning := fmt.Sprintf("%s. %s", definition.Meaning, item.Meaning)
		def := Definition{Meaning: appenedMeaning, Author: item.Author}
		if err := d.DB.Model(&definition).Updates(def).Error; err != nil {
			slog.Error("append definition error", "err", err)
			return err
		}
		return nil
	}
	return ErrLocked
}

// Lock a given definition (no more append or forget)
func (d *ZbotDatabaseSqlite) Lock(item Definition) error {
	definition, err := d.Get(item.Term, item.ChannelID)
	if err != nil {
		slog.Error("lock get error", "err", err)
		return err
	}

	if !definition.Locked.Bool {
		def := Definition{Locked: sql.NullBool{Bool: true, Valid: true}, LockedBy: sql.NullString{String: item.Author, Valid: true}}
		if err := d.DB.Model(&definition).Updates(def).Error; err != nil {
			slog.Error("lock update error", "err", err)
			return err
		}
		return nil
	}

	return fmt.Errorf("already locked by  %q", definition.LockedBy.String)
}

// Forget No implemented yet
func (d *ZbotDatabaseSqlite) Forget(item Definition) error {
	return nil
}
