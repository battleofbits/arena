CREATE TABLE fourup_moves (
    id serial,
    /* column is not an acceptable name */
    fourup_column int not null check (fourup_column >= 0),
    player int references players(id),
    move_id int not null,
    played timestamp not null default current_timestamp
);
