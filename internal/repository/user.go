package repository

import (
	"context"
	"gorm.io/gorm"
	"gormSession/internal/models"
	"gormSession/internal/query"
)

type User struct {
	q   *query.Query
	db  *gorm.DB
	ctx context.Context
}

func NewUser(ctx context.Context, db *gorm.DB) *User {
	return &User{
		ctx: ctx,
		db:  db,
		q:   query.Use(db),
	}
}

func (u *User) Create(tx *query.Query, m *models.User) error {
	q := u.q.User
	if tx != nil {
		q = tx.User
	}

	return q.WithContext(u.ctx).Create(m)
}

func (u *User) FindByName(tx *query.Query, name string) (*models.User, error) {
	q := u.q.User
	if tx != nil {
		q = tx.User
	}

	return q.WithContext(u.ctx).Where(q.Username.Eq(name)).First()
}
