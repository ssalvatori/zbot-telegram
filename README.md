# zbot-telegram

[![Go CICD](https://github.com/ssalvatori/zbot-telegram/actions/workflows/ci.yml/badge.svg)](https://github.com/ssalvatori/zbot-telegram/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/ssalvatori/zbot-telegram/graph/badge.svg?token=VWqgVPZrGy)](https://codecov.io/github/ssalvatori/zbot-telegram)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssalvatori/zbot-telegram)](https://goreportcard.com/report/github.com/ssalvatori/zbot-telegram)

A Telegram bot that serves as a collaborative knowledge base for chat groups. Users can teach the bot terms and definitions, search stored knowledge, and extend functionality through external modules.

## Features

- **Learn module** — Collaboratively build a knowledge database per chat group ([details](doc/learn.md))
- **External modules** — Extend the bot with custom scripts and binaries
- **User permissions** — Level-based access control for commands
- **Spam filtering** — Ignore system with configurable duration and automatic expiration
- **Per-chat isolation** — Each Telegram group has its own separate knowledge base
- **SQLite storage** — Simple, file-based persistence with automatic schema migrations

## Requirements

- Go 1.26+ (for building from source)
- A Telegram Bot API token from [BotFather](https://core.telegram.org/bots#botfather)

## Quick Start

### From Source

```bash
git clone https://github.com/ssalvatori/zbot-telegram.git
cd zbot-telegram
make build
```

### Configuration

Create a configuration file (default: `./zbot.conf`):

```yaml
zbot:
  token: <TELEGRAM_TOKEN>
  ignore_duration: 300
  ignore: true
  level: false
db:
  engine: sqlite
  file: ./zbot.db
commands:
  disabled:
    - ignore
    - level
    - forget
modules:
  path: ./modules/
  list: []
```

Set the config file path (optional if using `./zbot.conf`):

```bash
export ZBOT_CONFIG_FILE=/path/to/zbot.conf
./zbot-telegram
```

## Commands

### Learn Module

Commands for managing the knowledge base. Use `!` prefix for actions and `?` for retrieval.

| Command | Syntax | Description |
|---------|--------|-------------|
| learn | `!learn <term> <meaning>` | Store a new definition (auto-increments duplicates: term, term1, term2…) |
| get | `?<term>` | Retrieve the meaning of a term |
| who | `!who <term>` | Show metadata: author, date, hit count |
| append | `!append <term> <text>` | Append text to an existing term's meaning |
| find | `!find <text>` | Search inside meanings (limit: 10 results) |
| search | `!search <pattern>` | Search term names by pattern (limit: 10 results) |
| rand | `!rand [n]` | Get random terms (default: 1, max: 100) |
| top | `!top [n]` | Most retrieved terms by hit count (default: 10, max: 100) |
| last | `!last` | Show last 10 terms added |
| stats | `!stats` | Total number of definitions in current chat |
| lock | `!lock <term>` | Prevent a term from being modified (level 1000) |
| forget | `!forget <term>` | Delete a term (level 1000) |

### Utility Commands

| Command | Syntax | Description |
|---------|--------|-------------|
| ping | `!ping` | Health check — responds with "pong!!" |
| version | `!version` | Show bot version, git hash, and build time |
| ignore | `!ignore [list\|add\|help]` | Manage the spam filter ignore list (level 100) |
| level | `!level [add\|del] <user> <level>` | Manage user permission levels |

### External Modules

External modules are invoked with the `/` prefix:

```
/crypto btc
/test hello world
```

The bot executes the configured script passing: `<username> <user_level> <chat_name> <args>`

## Configuration Reference

See [doc/configuration.md](doc/configuration.md) for full details.

```yaml
zbot:
  token: <TELEGRAM_TOKEN>       # Required: Telegram Bot API token
  ignore_duration: 300          # Seconds to ignore a spamming user (default: 300)
  ignore: true                  # Enable spam filtering (default: false)
  level: false                  # Enable user permission levels (default: false)

db:
  engine: sqlite                # Database engine (sqlite)
  file: ./zbot.db               # Path to SQLite database file

commands:
  learn:
    disabled:                   # Channels where learn commands are disabled
      - channel_name
  disabled:                     # Globally disabled commands
    - ignore
    - level
    - forget

modules:
  path: ./modules/              # Path to external module executables
  list:
    - key: crypto               # Command name (invoked as /crypto)
      file: crypto              # Executable filename in modules path
      description: "Crypto data"
```

## External Modules

External modules allow extending the bot with any executable (shell scripts, Python, Go binaries, etc.).

### Creating a Module

1. Place the executable in the configured `modules.path` directory
2. The bot calls it with: `./module <username> <level> <chat> <user_args>`
3. The module's stdout is sent back as the bot's reply

Example module (`modules/hello`):

```bash
#!/bin/bash
echo "Hello $1! You are level $2 in chat $3. You said: $4"
```

4. Register it in the configuration:

```yaml
modules:
  list:
    - key: hello
      file: hello
      description: "Greets the user"
```

## Docker

### Build

```bash
make build-docker
```

Or with version info:

```bash
docker buildx build \
  --build-arg VERSION=$(git describe --tags) \
  --build-arg GIT_HASH=$(git rev-parse --short HEAD) \
  --build-arg BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ') \
  -t zbot-telegram .
```

### Run

```bash
docker run --rm \
  -v /path/to/zbot.conf:/app/zbot.conf \
  -v /path/to/modules:/app/modules \
  -v /path/to/zbot.db:/app/zbot.db \
  -e ZBOT_CONFIG_FILE=/app/zbot.conf \
  zbot-telegram:latest
```

### Multi-Architecture

Build for both amd64 and arm64:

```bash
make build-docker-multiarch
```

## Database

The bot uses [GORM](https://gorm.io/) for database management. All schemas and migrations are handled automatically on startup — no manual setup required.

The SQLite database file is created at the path specified in the configuration.

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Coverage

```bash
make coverage
make coverage-html   # Opens HTML report in browser
```

### Makefile Targets

| Target | Description |
|--------|-------------|
| `build` | Compile binary with version info |
| `test` | Run all tests with race detector |
| `coverage` | Generate test coverage report |
| `coverage-html` | Open HTML coverage in browser |
| `build-docker` | Build Docker image for current platform |
| `build-docker-multiarch` | Build for linux/amd64 and linux/arm64 |
| `build-docker-push` | Build and push multi-arch image to registry |
| `release` | Release with goreleaser |
| `lint` | Run `go vet` |
| `clean` | Remove build artifacts |

### Systemd Service

For running as a system service on Linux, see [extras/README.md](extras/README.md).

## License

See repository for license details.
