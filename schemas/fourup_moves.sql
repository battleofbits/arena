CREATE TABLE fourup_moves (
    id serial,
    column int not null check (column >= 0),
    player int references players(id),
    move_id int not null,
    played timestamp not null default current_timestamp
);
