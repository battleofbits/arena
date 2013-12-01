.PHONY: database arena deps test format

DATABASE=arena
USER=postgres_arena

database:
	# this top one must run as default user to grant permissions
	psql -f schemas/user.sql -U postgres
	psql -l | grep $(DATABASE) || psql -f schemas/database.sql --echo-all -U postgres
	psql -f schemas/players.sql -d $(DATABASE) --username=$(USER) --echo-all
	psql -f schemas/games.sql -d $(DATABASE) -U $(USER)
	psql -f schemas/fourup_matches.sql -d $(DATABASE) -U $(USER)
	psql -f schemas/fourup_moves.sql -d $(DATABASE) -U $(USER)

clean:
	psql -f schemas/reset.sql -d $(DATABASE) -U $(USER) || true
	psql -f schemas/reset_database.sql -U postgres || true

deps:
	go get -d -v ./...

format: deps
	go fmt ./...

test: format
	go test ./...
