-- Cleaning
DROP TABLE IF EXISTS conservatories;
DROP INDEX IF EXISTS conservatories_id;
DROP SEQUENCE IF EXISTS conservatories_id_seq;

-- Conservatories
CREATE SEQUENCE conservatories_id_seq;

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

CREATE UNIQUE INDEX conservatories_id ON conservatories (id);
