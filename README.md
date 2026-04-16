# Task CLI (Go)

A small command-line task manager built in Go.

## Features

- Add tasks
- List tasks
- Mark tasks as done
- Delete tasks
- Persist tasks in `data/tasks.json`

## Run

From the project root:

```bash
go run ./cmd/taskcli --help
```

Examples:

```bash
go run ./cmd/taskcli add "Learn Go"
go run ./cmd/taskcli list
go run ./cmd/taskcli done 1
go run ./cmd/taskcli delete 1
```

## Test

```bash
go test ./...
go test -race ./...
go vet ./...
```

## Project Layout

- `cmd/taskcli` - CLI entrypoint
- `internal/task` - domain logic and service
- `internal/storage/jsonstore` - JSON file store implementation
- `data/tasks.json` - persisted task data
