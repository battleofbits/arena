CREATE TABLE fourup_matches (
    id serial,
    player_red int references players(id),
    player_black int references players(id),
    winner int references players(id),
    started timestamp,
    finished timestamp,
    board int[][]
);
