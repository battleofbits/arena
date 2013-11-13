CREATE TABLE games (
    id serial,
    name varchar not null check (name <> '')
);
