package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	log "github.com/sirupsen/logrus"
)

type ZbotSqlite3Database struct {
	Db   *sql.DB
	File string
}

func (d *ZbotSqlite3Database) GetConnectionInfo() string {
	return fmt.Sprintf("DB: %s", d.File)
}

func (d *ZbotSqlite3Database) UserIgnoreList() ([]UserIgnore, error) {
	log.Debug("Getting ignore list")
	statement := "SELECT username, since, until FROM ignore_list"
	stmt, err := d.Db.Prepare(statement)
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

	var users []UserIgnore
	var user UserIgnore
	for rows.Next() {
		err2 := rows.Scan(&user.Username, &user.Since, &user.Until)
		if err2 != nil {
			return nil, err2
		}
		users = append(users, user)
	}
	return users, nil
}

//Init start connection with sql database
func (d *ZbotSqlite3Database) Init() error {

	if _, err := os.Stat(d.File); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Sqlite file [%s] does not exist!", d.File))
	}

	log.Info("Connecting to database: " + d.File)
	db, err := sql.Open("sqlite3", "file:"+d.File+"?cache=shared&mode=rwc")
	if err != nil {
		log.Error(err)
		return err
	}
	if db == nil {
		log.Error(err)
		return errors.New("Error connecting")
	}
	d.Db = db

	return nil
}

func (d *ZbotSqlite3Database) Close() {
	log.Debug("Closing connection")
	d.Db.Close()
}

//Statistics get the number of definitions in the db for a given chat
func (d *ZbotSqlite3Database) Statistics(chat string) (string, error) {
	statement := "select count(*) as total from definitions where chat = ?"
	var totalCount string
	err := d.Db.QueryRow(statement, chat).Scan(&totalCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return totalCount, errors.New("No Rows found")
		}
		return totalCount, err
	}

	return totalCount, err
}

func (d *ZbotSqlite3Database) Top() ([]Definition, error) {

	statement := "SELECT term FROM definitions ORDER BY hits DESC LIMIT 10"
	rows, err := d.Db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Definition
	for rows.Next() {
		var key string
		err2 := rows.Scan(&key)
		if err2 != nil {
			return nil, err2
		}
		items = append(items, Definition{Term: key})
	}

	return items, nil
}
func (d *ZbotSqlite3Database) Rand() (Definition, error) {
	var def Definition

	statement := "SELECT term, meaning FROM definitions ORDER BY random() LIMIT 1"
	rows, err := d.Db.Query(statement)
	if err != nil {
		return def, err
	}
	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(&def.Term, &def.Meaning)
		if err2 != nil {
			return def, err2
		}
	}

	return def, nil

}

func (d *ZbotSqlite3Database) Last() (Definition, error) {
	var def Definition
	statement := "SELECT term, meaning FROM definitions ORDER BY id DESC LIMIT 1"

	err := d.Db.QueryRow(statement).Scan(&def.Term, &def.Meaning)
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

//Get fetch a definition for a given term and chat
func (d *ZbotSqlite3Database) Get(term string, chat string) (Definition, error) {
	var def Definition
	statement := "SELECT id, term, meaning, chat, author, date FROM definitions WHERE term = ? and chat = ? COLLATE NOCASE LIMIT 1"
	err := d.Db.QueryRow(statement, term, chat).Scan(&def.ID, &def.Term, &def.Meaning, &def.Chat, &def.Author, &def.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return Definition{Term: "", Meaning: ""}, nil
		}
		log.Fatal(err)
		return def, err
	}

	statement = "UPDATE definitions SET hits = hits + 1 WHERE id = ?"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return def, err
	}

	_, err = stmt.Exec(def.ID)
	if err != nil {
		return def, err
	}

	return def, nil
}

func (d *ZbotSqlite3Database) _set(term string, def Definition) (sql.Result, error) {
	statement := "INSERT INTO definitions (term, meaning, author, locked, active, date, hits, link) VALUES (?,?,?,?,?,?,?,?)"

	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	return stmt.Exec(term, def.Meaning, def.Author, 1, 1, def.Date, 0, 0)

}

//Set add new definition in to the db
func (d *ZbotSqlite3Database) Set(def Definition) (string, error) {
	count := 1
	term := def.Term
	log.Debug(def)
	for {
		_, err := d._set(term, def)
		if err != nil {
			log.Debug("SQL insert error: ", err.Error())
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				term = fmt.Sprintf("%s%d", def.Term, count)
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

func (d *ZbotSqlite3Database) Append(def Definition) error {
	statement := "UPDATE definitions SET meaning = meaning || ?, date = ?, author = ? WHERE term = ?"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(fmt.Sprintf(" %s", def.Meaning), def.Date, def.Author, def.Term)
	if err != nil {
		log.Error(err.Error())
	}
	return nil
}

//Find search for a list of term with a pattern in the meaning for a given chat
func (d *ZbotSqlite3Database) Find(criteria string, chat string) ([]Definition, error) {
	var items []Definition
	statement := "SELECT term FROM definitions WHERE chat = ? and meaning like ? ORDER BY random() COLLATE NOCASE LIMIT 20"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return items, err
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(chat, criteria)
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
		items = append(items, Definition{Term: result})
	}
	return items, nil
}

//Search look for a list of term with a pattern in it for given chat
func (d *ZbotSqlite3Database) Search(criteria string, chat string) ([]Definition, error) {
	statement := "SELECT term FROM definitions WHERE chat = ?  and term like ? ORDER BY random() COLLATE NOCASE LIMIT 10"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(chat, criteria)
	if err2 != nil {
		panic(err2)
	}
	defer rows.Close()

	var items []Definition
	var result string
	for rows.Next() {
		err2 := rows.Scan(&result)
		if err2 != nil {
			return nil, err2
		}
		items = append(items, Definition{Term: result})
	}
	return items, nil
}

func (d *ZbotSqlite3Database) UserLevel(username string) (string, error) {
	var level string
	statement := "SELECT level FROM users WHERE username = ? COLLATE NOCASE LIMIT 1"
	err := d.Db.QueryRow(statement, username).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			return "0", nil
		} else {
			return level, err
		}
	}

	return level, nil
}
func (d *ZbotSqlite3Database) UserIgnoreInsert(username string) error {
	statement := "INSERT INTO ignore_list (username, since, until) VALUES (?,?,?)"
	stmt, err := d.Db.Prepare(statement)

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

//UserCheckIgnore check if user there is any row in ignore_list table for some username
// and an until greater than the current time
func (d *ZbotSqlite3Database) UserCheckIgnore(username string) bool {
	ignored := false

	now := time.Now().Unix()

	var level string
	statement := "SELECT count(*) as total FROM ignore_list WHERE username = ? AND until >= ?"
	err := d.Db.QueryRow(statement, username, now).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			ignored = false
		} else {
			log.Error(err)
			return ignored
		}
	}
	levelInt, _ := strconv.Atoi(level)
	log.Debug("Ingored ", levelInt)
	if levelInt > 0 {
		ignored = true
	}

	return ignored
}
func (d *ZbotSqlite3Database) UserCleanupIgnorelist() error {
	for {
		log.Debug("Cleaning ignore list")
		now := time.Now().Unix()
		statement := "DELETE FROM ignore_list WHERE until <= ?"
		stmt, err := d.Db.Prepare(statement)
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

func (d *ZbotSqlite3Database) Lock(item Definition) error {

	statement := "UPDATE definitions SET locked = 1, locked_by = ? WHERE term = ?"

	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(item.Author, item.Term)
	if err != nil {
		return err
	}

	return nil

}

func (d *ZbotSqlite3Database) Forget(item Definition) error {
	statement := "DELETE definitions WHERE term = ? LIMIT 1"

	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(item.Term)
	if err != nil {
		return err
	}

	return nil
}
