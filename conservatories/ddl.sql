DROP DATABASE IF EXISTS charts;
CREATE DATABASE charts;

DROP USER IF EXISTS charts;
CREATE USER charts WITH PASSWORD 'password';

GRANT ALL PRIVILEGES ON DATABASE charts to charts;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO charts;

DROP SEQUENCE IF EXISTS conservatories_id_seq;
CREATE SEQUENCE conservatories_id_seq;

DROP TABLE IF EXISTS conservatories;
CREATE TABLE conservatories (
  id INTEGER DEFAULT nextval('conservatories_id_seq') NOT NULL,
  name TEXT NOT NULL,
  category TEXT NOT NULL,
  street TEXT NOT NULL,
  city TEXT NOT NULL,
  department INTEGER NOT NULL,
  zip TEXT NOT NULL,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  creation_date TIMESTAMP DEFAULT now()
);

DROP INDEX IF EXISTS conservatories_id;
CREATE UNIQUE INDEX conservatories_id ON conservatories (id);
