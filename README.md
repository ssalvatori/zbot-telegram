# zbot-telegram

[![Build Status](https://travis-ci.org/ssalvatori/zbot-telegram.svg?branch=master)](https://travis-ci.org/ssalvatori/zbot-telegram)
[![Coverage Status](https://coveralls.io/repos/github/ssalvatori/zbot-telegram/badge.svg?branch=master)](https://coveralls.io/github/ssalvatori/zbot-telegram?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssalvatori/zbot-telegram)](https://goreportcard.com/report/github.com/ssalvatori/zbot-telegram)

## Requirements

* You need to get an API TOKEN from [BotFather@Telegram](https://core.telegram.org/bots)

## Setup

You **must** set the environment variable **ZBOT_CONFIG_FILE** with the path to the configuration file
* ZBOT_CONFIG_FILE : Path to the configuration file (default ./zbot.conf) 
* ZBOT_LOG_LEVEL : Log verbosity (default info)

## Configuration File

```yaml
zbot:
  token: <TELEGRAM_TOKEN>
  ignore_duration: 300
  ignore: true
  level: false
db:
  engine: sqlite
  file: path_to_sqlite_file.db
commands:  
  learn:
    disabled:
      - zbot_dev
  disabled:
    - ignore
    - level
    - forget
modules:
  path: ./modules/
  list:
    - key: crypto
      file: cypto
      description: get some crypto data
    - key: test
      file: test
      description: test module
    - key: temp
      file: temp.sh
      description: get weather info
    - key: plex
      file: plex2.py
      description: get plext information
```

## Database Schemas

[GORM](https://gorm.io/index.html), will creat the necessary schemas and migrations

## Development

```bash
docker build -t zbot-telegram -f Dockerfile .
docker run --rm --name zbot-telegram -v ${PWD}/modules:/app/modules -v ${PWD}/zbot.conf:/app/zbot.conf -e ZBOT_CONFIG_FILE=/app/zbot.conf zbot-telegram:latest
```
