DO
$body$
BEGIN
    IF NOT EXISTS (
        SELECT *
        FROM   pg_catalog.pg_user
        /* Need to put postgres in quotes, or this fails. */
        WHERE  usename = 'postgres') THEN

        CREATE USER postgres;
    END IF;
END
$body$
;

CREATE DATABASE arena WITH owner=postgres template=template0 encoding='UTF8'; 
