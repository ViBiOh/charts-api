-- Cleaning
DROP TABLE IF EXISTS readings_tags;
DROP TABLE IF EXISTS readings;
DROP TABLE IF EXISTS tags;

DROP INDEX IF EXISTS readings_id;
DROP INDEX IF EXISTS readings_username;
DROP INDEX IF EXISTS tags_id;
DROP INDEX IF EXISTS tags_username;
DROP INDEX IF EXISTS tags_name;

DROP SEQUENCE IF EXISTS readings_id_seq;
DROP SEQUENCE IF EXISTS tags_id_seq;

-- Readings
CREATE SEQUENCE readings_id_seq;

CREATE TABLE readings (
  id INTEGER NOT NULL DEFAULT nextval('readings_id_seq'),
  username TEXT NOT NULL,
  url TEXT NOT NULL,
  public BOOLEAN NOT NULL DEFAULT FALSE,
  read BOOLEAN NOT NULL DEFAULT FALSE,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX readings_id ON readings (id);
CREATE INDEX readings_username ON readings (username);

-- Tags
CREATE SEQUENCE tags_id_seq;

CREATE TABLE tags (
  id INTEGER NOT NULL DEFAULT nextval('tags_id_seq'),
  username TEXT NOT NULL,
  name TEXT NOT NULL,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX tags_id ON tags (id);
CREATE INDEX tags_username ON tags (username);
CREATE INDEX tags_name ON tags (name);

-- Tags / Readings
CREATE TABLE readings_tags (
  readings_id INTEGER NOT NULL REFERENCES readings(id),
  tags_id INTEGER NOT NULL REFERENCES tags(id),
  creation_date TIMESTAMP DEFAULT now()
);
