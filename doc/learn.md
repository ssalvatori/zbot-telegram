# Learn Module

The learn module is the core knowledge base system. It allows users in a Telegram group to collaboratively store and retrieve terms and definitions. Each chat group has its own isolated database.

## Configuration

The learn module is enabled by default. It can be disabled per-channel in the [configuration file](./configuration.md):

```yaml
commands:
  learn:
    disabled:
      - channel_name
```

## Commands

### Learn a term

Store a new definition. If the term already exists, an auto-incremented suffix is added (term, term1, term2…).

```
!learn <term> <meaning>
```

**Example:**
```
!learn golang A statically typed, compiled language designed at Google
```

### Get a term

Retrieve the meaning of a stored term. Each retrieval increments the term's hit counter.

```
?<term>
```

**Example:**
```
?golang
→ golang: A statically typed, compiled language designed at Google
```

### Who defined a term

Show metadata about a term: author, creation date, and number of times it has been retrieved.

```
!who <term>
```

### Append to a term

Add additional text to an existing term's meaning.

```
!append <term> <text>
```

**Example:**
```
!append golang — also great for CLI tools
```

### Search terms by name

Find terms whose names match a pattern. Returns up to 10 results.

```
!search <pattern>
```

**Example:**
```
!search go*
→ golang, goroutine, gomod
```

### Find terms by meaning content

Search inside the meanings of all terms. Returns up to 10 results.

```
!find <text>
```

**Example:**
```
!find compiled language
→ golang, rust, c
```

### Random terms

Get one or more random terms from the database.

```
!rand [number]
```

- Default: 1
- Maximum: 100

### Top terms

Get the most retrieved terms (highest hit count).

```
!top [number]
```

- Default: 10
- Maximum: 100

### Last terms added

Show the 10 most recently added terms.

```
!last
```

### Statistics

Show the total number of definitions stored in the current chat.

```
!stats
```

### Lock a term

Prevent a term from being modified or deleted. Requires level 1000.

```
!lock <term>
```

### Forget (delete) a term

Remove a term from the database (soft delete). Requires level 1000.

```
!forget <term>
```

