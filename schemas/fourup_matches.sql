CREATE TABLE fourup_matches (
    id serial,
    player_ids int[] references players(id),
    winner int references players(id),
    started timestamp,
    finished timestamp,
    board int[][]
);
