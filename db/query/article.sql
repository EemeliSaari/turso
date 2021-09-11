
-- name: CreateArticle :execresult
INSERT INTO articles (
    title,
    description,
    link,
    published,
    guid,
    is_loaded,
    checksum
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetArticle :one
SELECT * FROM articles WHERE id = $1 LIMIT 1;
