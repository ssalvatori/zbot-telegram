# zbot-telegram

[![Build Status](https://travis-ci.org/ssalvatori/zbot-telegram.svg?branch=master)](https://travis-ci.org/ssalvatori/zbot-telegram)
[![Coverage Status](https://coveralls.io/repos/github/ssalvatori/zbot-telegram/badge.svg?branch=dev)](https://coveralls.io/github/ssalvatori/zbot-telegram?branch=dev)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssalvatori/zbot-telegram)](https://goreportcard.com/report/github.com/ssalvatori/zbot-telegram)
[![Build status](https://ci.appveyor.com/api/projects/status/sml2sck29885bh4g?svg=true)](https://ci.appveyor.com/project/ssalvatori/zbot-telegram)
[![Build Status](https://cloud.drone.io/api/badges/ssalvatori/zbot-telegram/status.svg)](https://cloud.drone.io/ssalvatori/zbot-telegram)


## Requirements

* You need to get an API TOKEN from [BotFather@Telegram](https://core.telegram.org/bots)

## Setup

You **must** set the **ZBOT_TOKEN** environment variable using the Telegram's API TOKEN
 
* ZBOT_TOKEN : Use to connect to telegram
* ZBOT_LOG_LEVEL : Use to set the log level the alternatives are debug,info, warn, error, panic
* ZBOT_MODULES_PATH : Path to the externals modules directory
* ZBOT_DISABLED_COMMANDS: Commands to be ignored (separated by a ,)
* ZBOT_DATABASE_TYPE: Database to be use (mysql or sqlite)

## Feature Flags

* ZBOT_FLAG_ACTIVATE_IGNORE
* ZBOT_FLAG_ACTIVATE_LEVELS

## SQLite Configurations

* ZBOT_SQLITE_DATABASE: Path to the SQLite file 

## MySQL Configuration

* ZBOT_MYSQL_USERNAME : Database username
* ZBOT_MYSQL_PASSWORD : Database password
* ZBOT_MYSQL_DATABASE : Database name
* ZBOT_MYSQL_HOSTNAME : Database hostname 
* ZBOT_MYSQL_PORT : Database port (default: 3306)
* ZBOT_MYSQL_PROTOCOL: Protocol (default: tcp)

## Database Schemas

### Definitions

```sql
CREATE TABLE `definitions` ( 
    id INTEGER PRIMARY KEY AUTO_INCREMENT, 
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
