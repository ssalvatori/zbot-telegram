package db

import (
	"database/sql"

	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" //sql driver for mysql
	log "github.com/sirupsen/logrus"
)

//ZbotMysqlDatabase principal struct
type ZbotMysqlDatabase struct {
	Db         *sql.DB
	Connection MysqlConnection
}

//MysqlConnection databse configuration used to generate DSN
type MysqlConnection struct {
	Username     string
	Password     string
	DatabaseName string
	HostName     string
	Protocol     string
	Port         int
}

//GetConnectionInfo create connection string
func (d *ZbotMysqlDatabase) GetConnectionInfo() string {
	connectionDSN := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", d.Connection.Username, d.Connection.Password, d.Connection.Protocol, d.Connection.HostName, d.Connection.Port, d.Connection.DatabaseName)
	log.Debug("Using DSN: " + connectionDSN)
	return connectionDSN
}

//Init start a connection to database
func (d *ZbotMysqlDatabase) Init() error {
	log.Debug("Connecting to mysql database")
	connectionData := d.GetConnectionInfo()
	connection, err := sql.Open("mysql", connectionData)
	if err != nil {
		log.Error(err)
		return err
	}
	d.Db = connection
	return nil
}

//Close connection to database
func (d *ZbotMysqlDatabase) Close() {
	log.Debug("Closing connection")
	d.Db.Close()
}

//UserIgnoreList get list of ignored users
func (d *ZbotMysqlDatabase) UserIgnoreList() ([]UserIgnore, error) {
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

//Statistics get the number of definitions in the db for a given chat
func (d *ZbotMysqlDatabase) Statistics(chat string) (string, error) {
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

//Top the 10 definitions with more hits
//TODO: add parameter to change the number of definitions
func (d *ZbotMysqlDatabase) Top() ([]Definition, error) {

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

//Rand get a random definition from the database
func (d *ZbotMysqlDatabase) Rand() (Definition, error) {
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

func (d *ZbotMysqlDatabase) Last() (Definition, error) {
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

//Get fetch meaning for a definition in a given chat
func (d *ZbotMysqlDatabase) Get(term string, chat string) (Definition, error) {
	var def Definition
	statement := "SELECT id, term, meaning, author, date FROM definitions WHERE term = ? and chat = ? COLLATE NOCASE LIMIT 1"
	err := d.Db.QueryRow(statement, term).Scan(&def.ID, &def.Term, &def.Meaning, &def.Author, &def.Date, &def.Chat)
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

func (d *ZbotMysqlDatabase) _set(term string, def Definition) (sql.Result, error) {
	statement := "INSERT INTO definitions (term, meaning, chat, author, locked, active, date, hits, link) VALUES (?,?,?,?,?,?,?,?,?)"

	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	return stmt.Exec(term, def.Meaning, def.Chat, def.Author, 1, 1, def.Date, 0, 0)

}

//Set save new term in db
func (d *ZbotMysqlDatabase) Set(def Definition) (string, error) {
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

//Append append text to an existing definition
func (d *ZbotMysqlDatabase) Append(def Definition) error {
	statement := "UPDATE definitions SET meaning = meaning || ?, date = ?, author = ?, chat = ? WHERE term = ? and chat = ?"
	stmt, err := d.Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(fmt.Sprintf(" %s", def.Meaning), def.Date, def.Author, def.Chat, def.Term, def.Chat)
	if err != nil {
		log.Error(err.Error())
	}
	return nil
}

//Find get list of term with some string in the meaning
func (d *ZbotMysqlDatabase) Find(criteria string, chat string) ([]Definition, error) {
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

//Search get list of terms by some string in the name
func (d *ZbotMysqlDatabase) Search(criteria string, chat string) ([]Definition, error) {
	statement := "SELECT term FROM definitions WHERE chat = ? and term like ? ORDER BY random() COLLATE NOCASE LIMIT 10"
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

func (d *ZbotMysqlDatabase) UserLevel(username string) (string, error) {
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
func (d *ZbotMysqlDatabase) UserIgnoreInsert(username string) error {
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
func (d *ZbotMysqlDatabase) UserCheckIgnore(username string) bool {
	ignored := false

	now := time.Now().Unix()

	var level int
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

	log.Debug("Ingored ", level)
	if level > 0 {
		ignored = true
	}

	return ignored
}
func (d *ZbotMysqlDatabase) UserCleanupIgnorelist() error {
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

func (d *ZbotMysqlDatabase) Lock(item Definition) error {

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

func (d *ZbotMysqlDatabase) Forget(item Definition) error {
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
