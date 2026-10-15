# Zbot Documentation

## Contents

- [Overview](#overview)
- [Administrators](#administrators)
- [Channels](#channels)
- [Configuration](#configuration)
- [Help](#help)
- [Learning](#learning)
- [Levels](#levels)
- [Private Messages](#private-messages)
- [External Modules](#external-modules)

## Overview

To preserve the bot's IRC terminology, a `channel` in zbot refers to a Telegram group. A Telegram channel is a broadcast-only space and is not what zbot means by a channel.

## Commands

### Commands in a channel

Commands in a zbot channel should start with `!`

### Commands in a private message

Commands in a zbot private message should start with `/`, so they have to be register with telegram.


## Administrators

Bot administrators are listed under `zbot.admins` in the configuration. They can:

- Add the bot to Telegram channels.
- Read log lines using the private `/logs` command.
- Grant or revoke user levels in a channel.
- Enable or disable learning and external modules in a channel.

## Channels

- Each channel has its own learning dictionary.
- The bot should ignore invitations to a channel unless they come from one of that channel's administrators.
- Bot administrators can enable or disable external modules for a channel.

## Configuration

The configuration file contains the bot settings, database connection, administrators, command restrictions, and external module definitions. For example:

```yaml
zbot:
  token: <TELEGRAM_TOKEN>
  ignore_duration: 300
  ignore: true
  level: false
  admins: # List of bot's administrators (Telegram usernames)
    - admin1
    - admin2
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

- Learn a term: `!learn [term] [meaning]`
- Get a term's meaning: `?[term]`
- Get information about a term: `!who [term]`
- Append a meaning to a term: `!append [term] [meaning]`
- Get a random term: `!rand`
- Get the last learned term: `!last`
- Search for similar terms: `!search *term*`
- Find terms whose meaning contains text: `!find *text*`
- Get the top `[number]` terms: `!top [number]`
- Get the total number of terms: `!stats`

If a term already exists, a number is appended to the new term as a suffix, such as `term1`, `term2`, and so on.

## Levels

Levels are scoped to a channel. Channel administrators can do everything and can assign levels to other users. Users with a level can grant other users levels lower than their own.

- `!level` lists your current level in the channel.
- `!level update <user> <level>` adds or updates a user's level.
- `!level remove <user>` removes a user's level.

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