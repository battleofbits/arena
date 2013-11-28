DO
$body$
BEGIN
    IF NOT EXISTS (
        SELECT *
        FROM   pg_catalog.pg_user
        /* Need to put postgres in quotes, or this fails. */
        WHERE  usename = 'postgres_arena') THEN

        CREATE USER postgres_arena;
        ALTER USER postgres_arena CREATEDB;
        ALTER USER postgres_arena SUPERUSER;
    END IF;
END
$body$
;

