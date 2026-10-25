# Zbot-Telegram Documentation

[![Go CICD](https://github.com/ssalvatori/zbot-telegram/actions/workflows/ci.yml/badge.svg)](https://github.com/ssalvatori/zbot-telegram/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/ssalvatori/zbot-telegram/graph/badge.svg?token=VWqgVPZrGy)](https://codecov.io/github/ssalvatori/zbot-telegram)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssalvatori/zbot-telegram)](https://goreportcard.com/report/github.com/ssalvatori/zbot-telegram)

## Contents

- [Overview](#overview)
- [Requirements](#requirements)
- [Setup](#setup)
- [Administrators](#administrators)
- [Channels](#channels)
- [Configuration](#configuration)
- [Help](#help)
- [Learning](#learning)
- [Levels](#levels)
- [Private Messages](#private-messages)
- [External Modules](#external-modules)

## Overview

**Zbot-Telegram** is the telegram version of the infamous zbot on IRCzone and IRC.CL back in the 2000s

It was ported into golang and now it also a place to test some AI development practices.

## Important consideration

To preserve the bot's IRC terminology, a `channel` in zbot refers to a Telegram's group. In Telegram a channel is a broadcast-only space.

## Requirements

* You need to get an API TOKEN from [BotFather@Telegram](https://core.telegram.org/bots)

## Setup

You **must** set the **ZBOT_CONFIG_FILE** environment variable with the path to the configuration file

* ZBOT_CONFIG_FILE : Path to the configuration file (default ./zbot.conf)
* ZBOT_LOG_LEVEL : Log level (info, debug, error; default info)

### Configuration File

```yaml
zbot:
  token: <TELEGRAM_TOKEN>
  ignore_duration: 300
  ignore: true
  level: false
db:
  engine: sqlite
  file: path_to_sqlite_file.db
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
```

## Commands

There are 2 type of commands in the bot, the comands sent in a channel (telegram group) or via private messages.

### Commands in a channel

Commands in a zbot channel should start with `!`

### Commands in a private message

Commands in a zbot private message should start with `/`, so they have to be register with telegram.


## Bot Administrators

Bot administrators are listed under `zbot.admins` in the configuration. They can:

- Add the bot to Telegram channels. #TODO
- Read log lines using the private `/logs` command. #TODO
- Grant or revoke user levels in a channel. #TODO
- Enable or disable learning and external modules in a channel. #TODO

## Channels

- Each channel has its own learning dictionary.
- The bot should ignore invitations to a channel unless they come from one of that channel's administrators. #TODO
- Bot administrators can enable or disable external modules for a channel. #TODO

## Configuration

The configuration file contains the bot settings, database connection, administrators, command restrictions, and external module definitions. For example:

```yaml
zbot:
  token: <TELEGRAM_TOKEN>
  ignore_duration: 300
  ignore: true
  level: false
  admins: # List of bot administrators' Telegram user IDs
    - 123456789
    - 987654321
db:
  engine: sqlite
  file: path_to_sqlite_file.db
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

The database settings shown above configure SQLite. `engine` selects the database engine and `file` gives the database file path.

## Help

- `!help` displays the help message.
- `!help <command>` displays the help message for a specific command.

Each command should provide a `printHelp` function with a clear, concise description of what the command does and how to use it.

## Learning

The learning module stores terms and their meanings in the database and retrieves them on request. It is enabled by default and can be disabled for a channel by a bot administrator through the configuration or private-message commands.

- Learn a term: `!learn <term> <meaning>`, If a term already exists, a number is appended to the new term as a suffix, such as `term1`, `term2`, and so on.
- Get a term's meaning: `?<term>`
- Get information about a term: `!who <term>`
- Append a meaning to a term: `!append <term> <meaning>`
- Get a random term: `!rand`
- Get the last learned term: `!last`
- Search for similar terms: `!search <*term*>`
- Find terms whose meaning contains text: `!find <*text*>`
- Get the top `[number]` terms: `!top <number>`
- Get the total number of terms: `!stats`

## Levels

Levels are scoped to a channel. Channel administrators can do everything and can assign levels to other users. Users with a level can grant other users levels lower than their own.

- `!level` lists your current level in the channel. #TODO
- `!level update <user> <level>` adds or updates a user's level. #TODO
- `!level remove <user>` removes a user's level. #TODO

Bot administrators can also manage levels through private messages:

- `/level add <channel> <user> <level-number>` grant <level-number> to <user> in the <channel>
- `/level del <channel> <user>` removes all the levels from <user> in the <channel>

## Private Messages

Private messages from users who are not bot administrators are not allowed. The bot responds to those users with a message explaining that private messages are not allowed.

Bot administrators can use private messages to manage channel levels and modules:

- `/configuraiton learn <enable|disable> <channel>` enable or disabble learn capabilities for a given channel
- `/configuration external <enable|disable> <channel>` enable or disable external modules for a given channel
- `/logs` prints the recent log lines.

The documentation describes `/logs` as printing the last 10 lines.

## External Modules

External modules run in a separate process. They can be written in any language that can be executed on the system and can read arguments passed by the bot. Bot administrators can disable external modules for a channel through private messages.

Module arguments are passed as follows:

- `$1`: username of the user who executed the command
- `$2`: channel name where the command was executed
- `$3`: user's level in the channel
- `$4` through `$n`: arguments passed to the command


## Development

All the important commands for development can be found in [Makefile](./Makefile)

## Database Schemas

[GORM](https://gorm.io/index.html), will creat all the schemas and migrations