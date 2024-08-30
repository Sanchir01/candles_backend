package storeCategory

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type CategoryPostgresStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CategoryPostgresStore {
	return &CategoryPostgresStore{
		db: db,
	}
}

func (db *CategoryPostgresStore) CreateCategory(ctx context.Context) (string, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var categories []dbCategory

}

type dbCategory struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Description string    `db:"description"`
	Version     uint      `db:"version"`
}
