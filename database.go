package main

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
	"time"
)

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

type sqlLite struct {
	db   *sql.DB
	file string
}


func (d *sqlLite) userIgnoreList() ([]userIgnore, error) {
	log.Debug("Getting ignore list")
	statement := "SELECT username, since, until FROM ignore_list"
	stmt, err := d.db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err2 := stmt.Query()
	if err2 != nil {
		panic(err2)
	}
	defer rows.Close()

	var users []userIgnore
	var user userIgnore
	for rows.Next() {
		err2 := rows.Scan(&user.username, &user.since, &user.until)
		if err2 != nil {
			return nil, err2
		}
		users = append(users, user)
	}
	return users, nil
}

func (d *sqlLite) init() error {
	log.Debug("Connecting to database")
	db, err := sql.Open("sqlite3", d.file)
	if err != nil {
		log.Error(err)
		return err
	}
	if db == nil {
		log.Error(err)
		return errors.New("Error connecting")
	}
	d.db = db

	return nil
}

func (d *sqlLite) close() {
	log.Debug("Closing conecction")
	d.db.Close()
}

func (d *sqlLite) statistics() (string, error) {
	statement := "select count(*) as total from definitions"
	var totalCount string
	err := d.db.QueryRow(statement).Scan(&totalCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return totalCount, errors.New("No Rows found")
		} else {
			return totalCount, err
		}
	}

	return totalCount, err
}

func (d *sqlLite) top() ([]definitionItem, error) {

	statement := "SELECT term FROM definitions ORDER BY hits DESC LIMIT 10"
	rows, err := d.db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []definitionItem
	for rows.Next() {
		var key string
		err2 := rows.Scan(&key)
		if err2 != nil {
			return nil, err2
		}
		items = append(items, definitionItem{term: key})
	}

	return items, nil
}
func (d *sqlLite) rand() (definitionItem, error) {
	var def definitionItem

	sql := "SELECT term, meaning FROM definitions ORDER BY random() LIMIT 1"
	rows, err := d.db.Query(sql)
	if err != nil {
		return def, err
	}
	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(&def.term, &def.meaning)
		if err2 != nil {
			return def, err2
		}
	}

	return def, nil

}

func (d *sqlLite) last() (definitionItem, error) {
	var def definitionItem
	statement := "SELECT term, meaning FROM definitions ORDER BY id DESC LIMIT 1"

	err := d.db.QueryRow(statement).Scan(&def.term, &def.meaning)
	if err != nil {
		if err == sql.ErrNoRows {
			return def, errors.New("Nothing defined")
		} else {
			log.Fatal(err)
			return def, err
		}
	}

	return def, nil
}
func (d *sqlLite) get(term string) (definitionItem, error) {
	var def definitionItem
	statement := "SELECT id, term, meaning FROM definitions WHERE term = ? COLLATE NOCASE LIMIT 1"
	err := d.db.QueryRow(statement, term).Scan(&def.id, &def.term, &def.meaning)
	if err != nil {
		if err == sql.ErrNoRows {
			return definitionItem{term: "", meaning: ""}, nil
		} else {
			log.Fatal(err)
			return def, err
		}
	}

	statement = "UPDATE definitions SET hits = hits + 1 WHERE id = ?"
	stmt, err := d.db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return def, err
	}

	_, err = stmt.Exec(def.id)
	if err != nil {
		return def, err
	}

	return def, nil
}

func (d *sqlLite) _set(term string, def definitionItem) (sql.Result, error) {
	statement := "INSERT INTO definitions (term, meaning, author, locked, active, date, hits, link) VALUES (?,?,?,?,?,?,?,?)"

	stmt, err := d.db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	return stmt.Exec(term, def.meaning, def.author, 1, 1, def.date, 0, 0)

}

func (d *sqlLite) set(def definitionItem) error {
	count := 1
	term := def.term
	for {
		_, err := d._set(term, def)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: definitions.value") {
				term = fmt.Sprintf("%s%d", def.term, count)
				count = count + 1
			} else {
				return err
			}

		} else {
			break
		}
	}
	return nil

}

func (d *sqlLite) find(criteria string) ([]definitionItem, error) {
	var items []definitionItem
	statement := "SELECT term FROM definitions WHERE meaning like ? ORDER BY random() COLLATE NOCASE LIMIT 20"
	stmt, err := d.db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return items, err
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(criteria)
	if err2 != nil {
		return items, err
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		err2 := rows.Scan(&result)
		if err2 != nil {
			return items, err2
		}
		items = append(items, definitionItem{term: result})
	}
	return items, nil
}
func (d *sqlLite) search(criteria string) ([]definitionItem, error) {
	statement := "SELECT term FROM definitions WHERE term like ? ORDER BY random() COLLATE NOCASE LIMIT 10"
	stmt, err := d.db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(criteria)
	if err2 != nil {
		panic(err2)
	}
	defer rows.Close()

	var items []definitionItem
	var result string
	for rows.Next() {
		err2 := rows.Scan(&result)
		if err2 != nil {
			return nil, err2
		}
		items = append(items, definitionItem{term: result})
	}
	return items, nil
}

func (d *sqlLite) userLevel(username string) (string, error) {
	var level string
	statement := "SELECT level FROM users WHERE username = ? COLLATE NOCASE LIMIT 1"
	err := d.db.QueryRow(statement, username).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			return "0", nil
		} else {
			return level, err
		}
	}

	return level, nil
}
func (d *sqlLite) userIgnoreInsert(username string) error {
	statement := "INSERT INTO ignore_list (username, since, until) VALUES (?,?,?)"
	stmt, err := d.db.Prepare(statement)

	if err != nil {
		return err
	}

	since := time.Now().Unix()
	tenMinutes := 10 * time.Minute
	until := time.Now().Add(tenMinutes).Unix()

	_, err = stmt.Exec(username, since, until)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func (d *sqlLite) userCheckIgnore(username string) (bool, error) {
	ignored := false

	now := time.Now().Unix()

	var level string
	statement := "SELECT count(*) as total FROM ignore_list WHERE username = ? AND until >= ?"
	err := d.db.QueryRow(statement, username, now).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			ignored = false
		} else {
			log.Fatal(err)
			return ignored, err

		}
	}
	levelInt, _ := strconv.Atoi(level)
	log.Debug("Ingored ", levelInt)
	if levelInt > 0 {
		ignored = true
	}

	return ignored, nil
}
func (d *sqlLite) userCleanIgnore() error {
	for {
		log.Debug("Cleaning ignore list")
		now := time.Now().Unix()
		statement := "DELETE FROM ignore_list WHERE until <= ?"
		stmt, err := d.db.Prepare(statement)
		_, err = stmt.Query(now)
		if err != nil {
			if err == sql.ErrNoRows {

			} else {
				log.Fatal(err)
				return err
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

