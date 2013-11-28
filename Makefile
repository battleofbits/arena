.PHONY: database

DATABASE=arena

database:
	psql -l | grep arena || psql -f schemas/database.sql
	psql -f schemas/players.sql -d $(DATABASE)
	psql -f schemas/games.sql -d $(DATABASE)
	psql -f schemas/fourup_matches.sql -d $(DATABASE)
	psql -f schemas/fourup_moves.sql -d $(DATABASE)

clean:
	psql -f schemas/reset.sql -d $(DATABASE)
	psql -f schemas/reset_database.sql
