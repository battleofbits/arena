CREATE TABLE players (
    username varchar not null,
    name varchar not null check (name <> ''),
    id serial
);
