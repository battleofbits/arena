CREATE TABLE players (
    /* Your friendly name, for display in a UI */
    username varchar not null,
    /* Your unique id name */
    name varchar not null check (name <> ''),
    id serial primary key
);
ALTER TABLE players ADD CONSTRAINT unique_name UNIQUE (name);
