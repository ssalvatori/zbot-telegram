# Zbot Configuration

```yaml
zbot:
  token: <TELEGRAM_TOKEN>
  ignore_duration: 300
  ignore: true
  level: false
db:
  engine: sqlite
  file: path_to_sqlite_file.db
webhook:
  disable: true
  port: 13371
  auth:
    - channel: channel1
      id: 1234
      token: <YOUR_SECURE_TOKEN>
    - channel: channel2
      token: <YOUR_SECURE_TOKEN>
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

## zbot

## db
```yaml
db:
  engine: sqlite
  file: path_to_sqlite_file.db
```
## webhook
```yaml
webhook:
  disable: bool // Enable or disable bot webhook (default: false)
  port: int // Webhook port (default: 11337)
  auth:
    - channel: string // Channel name (bot will overwrite it using chat_id information)
      id: int64 // Telegram Chat_ID (leave it empty and the bot will try to get it using channel name)
      token: string // Token to autenticate request, this should be unique per channel
```
## commands

## modules