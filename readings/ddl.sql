DROP DATABASE IF EXISTS readings;
DROP USER IF EXISTS readings;

CREATE DATABASE readings;
CREATE USER readings WITH PASSWORD 'MotDePasseBaseDeDonneesReadings1';

GRANT ALL PRIVILEGES ON DATABASE readings to readings;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO readings;

-- Cleaning
DROP TABLE IF EXISTS readings_tags;
DROP TABLE IF EXISTS readings;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS readings_id;
DROP INDEX IF EXISTS readings_user;
DROP INDEX IF EXISTS tags_id;
DROP INDEX IF EXISTS tags_name;
DROP INDEX IF EXISTS tags_user;
DROP INDEX IF EXISTS users_id;
DROP INDEX IF EXISTS users_id;
DROP INDEX IF EXISTS users_name;
DROP INDEX IF EXISTS users_name;

DROP SEQUENCE IF EXISTS readings_id_seq;
DROP SEQUENCE IF EXISTS tags_id_seq;
DROP SEQUENCE IF EXISTS users_id_seq;
DROP SEQUENCE IF EXISTS users_id_seq;

-- Users table
CREATE SEQUENCE users_id_seq;

CREATE TABLE users (
  id INTEGER NOT NULL DEFAULT nextval('users_id_seq'),
  name TEXT NOT NULL,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX users_id ON users (id);
CREATE INDEX users_name ON users (name);

-- Readings table
CREATE SEQUENCE readings_id_seq;

CREATE TABLE readings (
  id INTEGER NOT NULL DEFAULT nextval('readings_id_seq'),
  user_id INTEGER NOT NULL REFERENCES users(id),
  url TEXT NOT NULL,
  public BOOLEAN NOT NULL DEFAULT FALSE,
  read BOOLEAN NOT NULL DEFAULT FALSE,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX readings_id ON readings (id);

-- Tags table
CREATE SEQUENCE tags_id_seq;

CREATE TABLE tags (
  id INTEGER NOT NULL DEFAULT nextval('tags_id_seq'),
  user_id INTEGER NOT NULL REFERENCES users(id),
  name TEXT NOT NULL,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX tags_id ON tags (id);
CREATE INDEX tags_name ON tags (name);

-- Tags / Readings table
CREATE TABLE readings_tags (
  readings_id INTEGER NOT NULL REFERENCES readings(id),
  tags_id INTEGER NOT NULL REFERENCES tags(id),
  creation_date TIMESTAMP DEFAULT now()
);
