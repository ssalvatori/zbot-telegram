# zbot-telegram

[![Build Status](https://travis-ci.org/ssalvatori/zbot-telegram-go.svg?branch=master)](https://travis-ci.org/ssalvatori/zbot-telegram-go)
[![Coverage Status](https://coveralls.io/repos/github/ssalvatori/zbot-telegram-go/badge.svg)](https://coveralls.io/github/ssalvatori/zbot-telegram-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssalvatori/zbot-telegram-go)](https://goreportcard.com/report/github.com/ssalvatori/zbot-telegram-go)

## Requirements

* You need to get an API TOKEN from [BotFather@Telegram](https://core.telegram.org/bots)

## Setup

You **must** set the **ZBOT_TOKEN** environment variable using the Telegram's API TOKEN
 
* ZBOT_TOKEN : Use to connect to telegram
* ZBOT_LOG_LEVEL : Use to set the log level the alternatives are debug,info, warn, error, panic
* ZBOT_DATABASE_FILE : Path to the sqlite database "/path/to/file.sqlite"
* ZBOT_MODULES_PATH : Path to the externals modules directory

## Database Schemas

### Definitions

```sql
CREATE TABLE `definitions` ( 
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    term TEXT UNIQUE, 
    meaning TEXT NOT NULL, 
    author TEXT NOT NULL, 
    locked INTEGER DEFAULT 0, 
    active INTEGER DEFAULT 1, 
    date TEXT DEFAULT '0000-00-00', 
    hits INTEGER DEFAULT 0, 
    link INTEGER ,
    locked_by TEXT
)
```

### Users

```sql
CREATE TABLE `users` ( 
    `id` INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, 
    `username` TEXT NOT NULL, 
    `level` INTEGER DEFAULT 10 
)
```


You need at least one user in the database

```sql
INSERT INTO users VALUES (null, 'ssalvato', 1000)
```

### Ignore

```sql
CREATE TABLE `ignore_list` ( 
    `id` INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, 
    `username` TEXT NOT NULL, 
    `since` INTEGER NOT NULL,
    `until` INTEGER NOT NULL
)
```

### Migration (optional)

Ths is just for the migration of the oldest database schema

```sql
ALTER TABLE `ledger` RENAME TO `definitions`
```
