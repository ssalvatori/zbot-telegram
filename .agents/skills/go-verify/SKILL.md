---
name: go-verify
description: Run tests, linting, and formatting checks for this Go project. Use after modifying any .go file, before committing, or when asked to verify, test, or lint the code.
---

# Verify Go Changes

```bash
gofmt -l .
goimports -l .
```

Both commands should print nothing. If files are listed, fix them with `gofmt -w <file>` and `goimports -w <file>`.

## 2. Vet and lint

```bash
go vet ./...
golangci-lint run ./...
```

Do not silence linter warnings with `//nolint` unless there is a clear reason, and always include a comment explaining why.

## 3. Test

```bash
go test ./... -race -count=1
```

## 4. Tidy modules

If you added or removed imports, run:

```bash
go mod tidy
```

