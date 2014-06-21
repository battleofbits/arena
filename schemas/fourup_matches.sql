CREATE TABLE fourup_matches (
    id serial primary key,
    player_red int references players(id),
    player_black int references players(id),
    winner int references players(id),
    started timestamp,
    finished timestamp,
	/* unfortunately the go postgres driver doesn't support arrays, so this
	 * gets serialized to json before getting dumped to the DB. there are
	 * better optimization formats for storing the data besides JSON but in
	 * this case let's optimize for programmer productivity and figure it out
	 * if/when this actually becomes a thing people use. */
    board varchar
);
