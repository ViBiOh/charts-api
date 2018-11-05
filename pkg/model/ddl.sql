-- clean
DROP TABLE IF EXISTS reading_tag;
DROP TABLE IF EXISTS reading;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS "user";

DROP INDEX IF EXISTS user_id;
DROP INDEX IF EXISTS user_login;
DROP INDEX IF EXISTS reading_id;
DROP INDEX IF EXISTS reading_user;
DROP INDEX IF EXISTS tag_id;
DROP INDEX IF EXISTS tag_user;
DROP INDEX IF EXISTS tag_name;

-- user
CREATE TABLE "user" (
  id TEXT NOT NULL,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX user_id ON "user" (id);
CREATE UNIQUE INDEX user_login ON "user" (username);

-- reading
CREATE TABLE reading (
  id TEXT NOT NULL,
  user_id TEXT NOT NULL REFERENCES "user"(id),
  url TEXT NOT NULL,
  read BOOLEAN NOT NULL DEFAULT FALSE,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX reading_id ON reading (id);
CREATE INDEX reading_user ON reading (user_id);

-- tag
CREATE TABLE tag (
  id TEXT NOT NULL,
  user_id TEXT NOT NULL REFERENCES "user"(id),
  name TEXT NOT NULL,
  creation_date TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX tag_id ON tag (id);
CREATE INDEX tag_user ON tag (user_id);
CREATE INDEX tag_name ON tag (name);

-- tag / reading association
CREATE TABLE reading_tag (
  reading_id TEXT NOT NULL REFERENCES reading(id),
  tag_id TEXT NOT NULL REFERENCES tag(id),
  creation_date TIMESTAMP DEFAULT now()
);
