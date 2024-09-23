ALTER TABLE users
    ADD COLUMN username TEXT NOT NULL UNIQUE,
    ADD COLUMN location TEXT,
    ADD COLUMN birth_date DATE,
    ADD COLUMN sex TEXT;

