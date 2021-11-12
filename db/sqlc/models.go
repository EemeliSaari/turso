// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
)

type Article struct {
	ID               int32
	Title            sql.NullString
	ShortDescription sql.NullString
	Link             sql.NullString
	Published        interface{}
	ArticleGuid      sql.NullString
	IsLoaded         sql.NullBool
	ArticleChecksum  sql.NullString
}