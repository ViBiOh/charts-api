-- clean
DROP TABLE IF EXISTS reading_tag;
DROP TABLE IF EXISTS reading;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS "user";

DROP SEQUENCE IF EXISTS user_seq;
DROP SEQUENCE IF EXISTS reading_seq;
DROP SEQUENCE IF EXISTS tag_seq;

DROP INDEX IF EXISTS user_id;
DROP INDEX IF EXISTS user_login;
DROP INDEX IF EXISTS reading_id;
DROP INDEX IF EXISTS reading_user;
DROP INDEX IF EXISTS tag_id;
DROP INDEX IF EXISTS tag_user;
DROP INDEX IF EXISTS tag_name;

-- user
CREATE SEQUENCE user_seq;
CREATE TABLE "user" (
  id BIGINT NOT NULL DEFAULT nextval('user_seq'),
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  creation_date TIMESTAMP WITH TIME ZONE DEFAULT now()
);
ALTER SEQUENCE user_seq OWNED BY "user".id;

CREATE UNIQUE INDEX user_id ON "user" (id);
CREATE UNIQUE INDEX user_login ON "user" (username);

-- reading
CREATE SEQUENCE reading_seq;
CREATE TABLE reading (
  id BIGINT NOT NULL DEFAULT nextval('reading_seq'),
  user_id BIGINT NOT NULL REFERENCES "user"(id),
  url TEXT NOT NULL,
  read BOOLEAN NOT NULL DEFAULT FALSE,
  creation_date TIMESTAMP WITH TIME ZONE DEFAULT now()
);
ALTER SEQUENCE reading_seq OWNED BY reading.id;

CREATE UNIQUE INDEX reading_id ON reading (id);
CREATE INDEX reading_user ON reading (user_id);

-- tag
CREATE SEQUENCE tag_seq;
CREATE TABLE tag (
  id BIGINT NOT NULL DEFAULT nextval('tag_seq'),
  user_id BIGINT NOT NULL REFERENCES "user"(id),
  name TEXT NOT NULL,
  creation_date TIMESTAMP WITH TIME ZONE DEFAULT now()
);
ALTER SEQUENCE tag_seq OWNED BY tag.id;

CREATE UNIQUE INDEX tag_id ON tag (id);
CREATE INDEX tag_user ON tag (user_id);
CREATE INDEX tag_name ON tag (name);

-- tag / reading association
CREATE TABLE reading_tag (
  reading_id BIGINT NOT NULL REFERENCES reading(id),
  tag_id BIGINT NOT NULL REFERENCES tag(id),
  creation_date TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX reading_tag_reading_id ON reading_tag (reading_id);
CREATE INDEX reading_tag_tag_id ON reading_tag (tag_id);
