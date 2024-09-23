default:
    @just --list

# build soteria binary
build:
    go build -o soteria ./cmd/soteria

# update go packages
update:
    @cd ./cmd/soteria && go get -u

# set up the dev environment with docker-compose
dev cmd *flags:
    #!/usr/bin/env bash
    set -euxo pipefail
    if [ {{ cmd }} = 'down' ]; then
      docker compose -f ./docker-compose.yml down
      docker compose -f ./docker-compose.yml rm
    elif [ {{ cmd }} = 'up' ]; then
      docker compose -f ./docker-compose.yml up --wait -d {{ flags }}
    else
      docker compose -f ./docker-compose.yml {{ cmd }} {{ flags }}
    fi

# run tests in the dev environment
test:
    go test -v ./... -covermode=atomic -coverprofile=coverage.out

# run golangci-lint
lint:
    golangci-lint run -c .golangci.yml
