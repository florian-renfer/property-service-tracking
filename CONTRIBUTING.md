# Contributing

## Prerequisites

1. Go version from `go.mod`
2. Docker + Docker Compose

## Local workflow

1. `cp .env.example .env`
2. `make tools`
3. `make fmt`
4. `make lint`
5. `make test`
6. `make vet`

## Pull requests

1. Keep changes scoped and small.
2. Ensure local checks pass before opening PR.
3. CI must be green before merge.
