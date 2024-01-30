CREATE TABLE IF NOT EXISTS urls(
    id INTEGER PRIMARY KEY,
    alias TEXT NOT NULL unique,
    url TEXT NOT null

);
    CREATE INDEX IF NOT EXISTS idx_alias on url(alias);