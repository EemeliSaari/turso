CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title TEXT,
    short_description TEXT,
    link TEXT,
    published TIMESTAMP,
    article_guid TEXT,
    is_LOADED BOOLEAN,
    article_checksum TEXT
);