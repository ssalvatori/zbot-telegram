# zbot-telegram

## requirements

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
    link INTEGER 
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