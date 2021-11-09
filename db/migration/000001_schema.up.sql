CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title TEXT,
    short_description TEXT,
    link TEXT,
    published DATETIME,
    article_guid TEXT,
    is_LOADED BOOL,
    article_checksum TEXT
);
