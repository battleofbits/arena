CREATE TABLE players (
    username varchar not null,
    name varchar not null check (name <> ''),
    id serial primary key
);
ALTER TABLE players ADD CONSTRAINT unique_name UNIQUE (name);
