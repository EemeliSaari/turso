CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    link TEXT,
    published DATETIME,
    guid TEXT,
    is_LOADED BOOL,
    checksum TEXT
);
