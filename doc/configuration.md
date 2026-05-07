# Zbot Configuration

The bot is configured via a YAML file. By default it reads `./zbot.conf`, override with the `ZBOT_CONFIG_FILE` environment variable.

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
      file: crypto
      description: get some crypto data
```

## Sections

### zbot

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `token` | string | — | **Required.** Telegram Bot API token from BotFather |
| `ignore_duration` | int | 300 | Duration in seconds to ignore a user after spam detection |
| `ignore` | bool | false | Enable the spam filter / ignore system |
| `level` | bool | false | Enable user permission level enforcement |

### db

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `engine` | string | sqlite | Database engine (`sqlite`) |
| `file` | string | — | Path to the SQLite database file |

The database and all tables are created automatically on first startup via GORM migrations.

### commands

```yaml
commands:
  learn:
    disabled:           # List of chat names where learn commands are disabled
      - channel_name
  disabled:             # List of globally disabled commands
    - ignore
    - level
    - forget
```

Disabling a command removes it from the bot entirely — users will receive no response.

### modules

```yaml
modules:
  path: ./modules/      # Directory containing module executables
  list:
    - key: crypto       # Command name (invoked as /crypto)
      file: crypto      # Filename of the executable in modules path
      description: "Get cryptocurrency data"
```

Each module is called with arguments: `<username> <user_level> <chat_name> <user_args>`

The module's stdout is returned as the bot's reply message.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ZBOT_CONFIG_FILE` | Path to the configuration file (default: `./zbot.conf`) |