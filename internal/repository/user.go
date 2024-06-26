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

func NewUser(ctx context.Context) *User {
	return &User{
		ctx: ctx,
		//db:  db,
		q: GetQuery(),
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

func (u *User) FindSubQuery(tx *query.Query, name string) (*models.User, error) {
	q := u.q.User
	if tx != nil {
		q = tx.User
	}

	subQuery := q.Select(q.ID).WithContext(u.ctx).Where(q.Username.Eq(name))
	return q.WithContext(u.ctx).Where(q.Columns(q.ID).In(subQuery)).First()
}
