.PHONY: database arena

DATABASE=arena
USER=postgres_arena

database:
	# this top one must run as default user to grant permissions
	psql -f schemas/user.sql 
	psql -l | grep $(DATABASE) || psql -f schemas/database.sql --echo-all
	psql -f schemas/players.sql -d $(DATABASE) --username=$(USER) --echo-all
	psql -f schemas/games.sql -d $(DATABASE) -U $(USER)
	psql -f schemas/fourup_matches.sql -d $(DATABASE) -U $(USER)
	psql -f schemas/fourup_moves.sql -d $(DATABASE) -U $(USER)

clean:
	psql -f schemas/reset.sql -d $(DATABASE) -U $(USER) || true
	psql -f schemas/reset_database.sql || true

format:
	go fmt ./...

test: format
	go test ./...
