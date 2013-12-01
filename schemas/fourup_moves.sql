CREATE TABLE fourup_moves (
    id serial,
    /* column is not an acceptable Postgres name */
    fourup_column int not null check (fourup_column >= 0),
    player int references players(id),
    /* The move # in the game */
    move_number int not null,
    match_id int not null references fourup_matches(id),
    played timestamp not null default current_timestamp
);
