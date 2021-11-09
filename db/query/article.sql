-- name: CreateArticle :execresult
INSERT INTO articles (
    title,
    short_description,
    link,
    published,
    article_guid,
    is_loaded,
    article_checksum
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetArticle :one
SELECT * FROM articles WHERE id = $1 LIMIT 1;

-- name: ListArticles :many
SELECT * FROM articles ORDER BY name;
