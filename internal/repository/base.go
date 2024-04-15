package repository

import (
	"gormSession/internal/db"
	"gormSession/internal/query"
)

func GetQuery() *query.Query {
	return query.Use(db.GetDb())
}
