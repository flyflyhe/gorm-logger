package repository

import (
	"gorm.io/gorm"
	"gormSession/internal/query"
)

func GetQuery(db *gorm.DB) *query.Query {
	return query.Use(db)
}
