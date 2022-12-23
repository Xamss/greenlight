ALTER TABLE movies DROP COLUMN version;
ALTER TABLE movies ADD COLUMN version uuid NOT NULL default gen_random_uuid();