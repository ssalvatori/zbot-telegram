package main

import (
	"database/sql"
	"strings"
	"fmt"
	"time"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/Sirupsen/logrus"
)

type definitionItem struct {
	term	string
	meaning	string
	author string
	date string
}

func InitDB(filepath string) *sql.DB {
	log.Debug("Connecting to database")
	db, err := sql.Open("sqlite3", filepath)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func getStats(db *sql.DB) string {
	statement := "select count(*) as total from definitions"
	var totalCount string;
	err := db.QueryRow(statement).Scan(&totalCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return "error"
		} else {
			log.Fatal(err)
		}
	}

	return totalCount
}

func getTop(db *sql.DB) []string {

	statement := "SELECT term FROM definitions ORDER BY hits DESC LIMIT 10"
	rows, err := db.Query(statement)
	if err != nil { panic(err) }
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		err2 := rows.Scan(&key)
		if err2 != nil {
			panic(err2)
		}
		keys = append(keys, key)
	}

	return keys

}

func getRand(db *sql.DB) definitionItem {

	sql := "SELECT term, meaning FROM definitions ORDER BY random() LIMIT 1"
	rows, err := db.Query(sql)
	if err != nil { panic(err) }
	defer rows.Close()

	var def definitionItem
	for rows.Next() {
		err2 := rows.Scan(&def.term, &def.meaning)
		if err2 != nil {
			panic(err2)
		}
	}

	return def

}

func getLast(db *sql.DB) definitionItem {
	var def definitionItem
	statement := "SELECT term, meaning FROM definitions ORDER BY id DESC LIMIT 1"

	err := db.QueryRow(statement).Scan(&def.term, &def.meaning)
	if err != nil {
		if err == sql.ErrNoRows {
			return definitionItem{term: "", meaning: ""}
		} else {
			log.Fatal(err)
		}
	}

	return def
}

func getDefinition(db *sql.DB, term string) definitionItem {
	var def definitionItem
	statement := "SELECT term, meaning FROM definitions WHERE term = ? COLLATE NOCASE LIMIT 1"
	err := db.QueryRow(statement, term).Scan(&def.term, &def.meaning)
	if err != nil {
		if err == sql.ErrNoRows {
			return definitionItem{term: "", meaning: ""}
		} else {
			log.Fatal(err)
		}
	}

	return def
}

func setDefinition(db *sql.DB, def definitionItem) (string, string)  {
	count := 1
	term := def.term
	for {
		_, err := insertDef(db, term, def)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: definitions.value") {
				term = fmt.Sprintf("%s%d", def.term, count)
				count = count + 1
			}
		} else {
			break
		}
	}
	return term, def.meaning

}

func insertDef(db *sql.DB, term string, def definitionItem) (sql.Result, error) {
	statement := "INSERT INTO definitions (term, meaning, author, locked, active, date, hits, link) VALUES (?,?,?,?,?,?,?,?)"

	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	return stmt.Exec(term, def.meaning, def.author, 1, 1, def.date, 0, 0)

}

func findDef(db *sql.DB, criteria string) []string {
	statement := "SELECT term FROM definitions WHERE meaning like ? ORDER BY random() COLLATE NOCASE LIMIT 20"
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(criteria)
	if err2 != nil { panic(err2) }
	defer rows.Close()

	var results []string
	for rows.Next() {
		result := ""
		err2 := rows.Scan(&result)
		if err2 != nil { panic(err2) }
		results = append(results, result)
	}
	return results
}

func searchDef(db *sql.DB, criteria string) []string {
	statement := "SELECT term FROM definitions WHERE term like ? ORDER BY random() COLLATE NOCASE LIMIT 10"
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err2 := stmt.Query(criteria)
	if err2 != nil { panic(err2) }
	defer rows.Close()

	var results []string
	for rows.Next() {
		result := ""
		err2 := rows.Scan(&result)
		if err2 != nil { panic(err2) }
		results = append(results, result)
	}
	return results
}

func getLevel(db *sql.DB, username string) string {
	var level string
	statement := "SELECT level FROM users WHERE username = ? COLLATE NOCASE LIMIT 1"
	err := db.QueryRow(statement, username).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			return "0"
		} else {
			log.Fatal(err)
		}
	}

	return level
}

func insertIgnoreUser(db *sql.DB, username string) {
	statement := "INSERT INTO ignore_list (username, since, until) VALUES (?,?,?)"
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}

	since := time.Now().Unix()
	tenMinutes := 10 * time.Minute
	until := time.Now().Add(tenMinutes).Unix()

	_, err = stmt.Exec(username, since, until)
	if err != nil {
		log.Fatal(err)
	}
}

func checkIgnoreList(db *sql.DB, username string) bool {
	ignored := false

	now := time.Now().Unix()

	var level string
	statement := "SELECT count(*) as total FROM ignore_list WHERE username = ? AND until >= ?"
	err := db.QueryRow(statement, username, now).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			ignored = false
		} else {
			log.Fatal(err)
		}
	}
	levelInt, _ := strconv.Atoi(level)
	log.Debug("Ingored ", levelInt)
	if ( levelInt > 0) {
		ignored = true
	}

	return ignored
}

func cleanIgnoreList() {
	for {
		log.Debug("Cleaning ignore list")
		now := time.Now().Unix()
		statement := "DELETE FROM ignore_list WHERE until <= ?"
		stmt, err := db.Prepare(statement)
		_, err = stmt.Query(now)
		if err != nil {
			if err == sql.ErrNoRows {

			} else {
				log.Fatal(err)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}