.PHONY: database arena deps test format

export GOPATH:=$(shell pwd)

DATABASE=arena
USER=postgres_arena
DEFAULT_USER=""

ifeq ($(TRAVIS), true)
DEFAULT_USER="-U postgres"
endif

database:
	# this top one must run as default user to grant permissions
	psql -f schemas/user.sql $(DEFAULT_USER)
	psql -l | grep $(DATABASE) || psql -f schemas/database.sql --echo-all $(DEFAULT_USER)
	psql -f schemas/players.sql -d $(DATABASE) --username=$(USER) --echo-all
	psql -f schemas/games.sql -d $(DATABASE) -U $(USER)
	psql -f schemas/fourup_matches.sql -d $(DATABASE) -U $(USER)
	psql -f schemas/fourup_moves.sql -d $(DATABASE) -U $(USER)

clean:
	psql -f schemas/reset.sql -d $(DATABASE) -U $(USER) || true
	psql -f schemas/reset_database.sql $(DEFAULT_USER) || true

deps:
	go get -d -v ./...
	go test -i ./...

format: deps
	go fmt ./...

test: format
	go test ./...

serve: format
	go build ./arena ./server
	go install ./server
	server

arena:
	go install arena

test:
	go test arena

mock: format
	go build ./arena ./mock
	go run mock/mock.go
