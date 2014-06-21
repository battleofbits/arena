CREATE TABLE players (
    id serial primary key,
    /* Your friendly name, for display in a UI */
    username varchar not null,
    /* Your unique id name */
    name varchar not null check (name <> ''),
    match_url varchar not null,
    invite_url varchar not null
);
ALTER TABLE players ADD CONSTRAINT unique_name UNIQUE (name);
