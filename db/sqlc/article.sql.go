// Code generated by sqlc. DO NOT EDIT.
// source: article.sql

package db

import (
	"context"
	"database/sql"
)

const createArticle = `-- name: CreateArticle :execresult
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
) RETURNING id, title, short_description, link, published, article_guid, is_loaded, article_checksum
`

type CreateArticleParams struct {
	Title            sql.NullString
	ShortDescription sql.NullString
	Link             sql.NullString
	Published        interface{}
	ArticleGuid      sql.NullString
	IsLoaded         sql.NullBool
	ArticleChecksum  sql.NullString
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createArticle,
		arg.Title,
		arg.ShortDescription,
		arg.Link,
		arg.Published,
		arg.ArticleGuid,
		arg.IsLoaded,
		arg.ArticleChecksum,
	)
}

const getArticle = `-- name: GetArticle :one
SELECT id, title, short_description, link, published, article_guid, is_loaded, article_checksum FROM articles WHERE id = $1 LIMIT 1
`

func (q *Queries) GetArticle(ctx context.Context, id int32) (Article, error) {
	row := q.db.QueryRowContext(ctx, getArticle, id)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.ShortDescription,
		&i.Link,
		&i.Published,
		&i.ArticleGuid,
		&i.IsLoaded,
		&i.ArticleChecksum,
	)
	return i, err
}

const listArticles = `-- name: ListArticles :many
SELECT id, title, short_description, link, published, article_guid, is_loaded, article_checksum FROM articles ORDER BY name
`

func (q *Queries) ListArticles(ctx context.Context) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, listArticles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.ShortDescription,
			&i.Link,
			&i.Published,
			&i.ArticleGuid,
			&i.IsLoaded,
			&i.ArticleChecksum,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
