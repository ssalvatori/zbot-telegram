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
## commands

## modules