package pgstoreCategory

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"log/slog"
	"time"
)

type CategoryPostgresStore struct {
	db *sqlx.DB
	lg *slog.Logger
}

func New(db *sqlx.DB) *CategoryPostgresStore {
	return &CategoryPostgresStore{
		db: db,
	}
}
func (db *CategoryPostgresStore) CategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	var category dbCategory
	if err := conn.GetContext(ctx, &category, "SELECT * FROM category WHERE slug = $1", slug); err != nil {
		return nil, err
	}
	return (*model.Category)(&category), nil

}
func (db *CategoryPostgresStore) AllCategories(ctx context.Context) ([]model.Category, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	var categories []dbCategory

	if err := conn.SelectContext(ctx, &categories, "SELECT * FROM category"); err != nil {
		return nil, err
	}
	return lo.Map(categories, func(category dbCategory, _ int) model.Category { return model.Category(category) }), nil
}

func (db *CategoryPostgresStore) CreateCategory(ctx context.Context, name, slug string) (uuid.UUID, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	defer conn.Close()

	var id uuid.UUID

	row := conn.QueryRowContext(ctx,
		"INSERT INTO category(name, slug) VALUES($1, $2) RETURNING id",
		name, slug,
	)
	if err := row.Err(); err != nil {
		db.lg.Error("create category error", err.Error())
		return uuid.New(), err
	}
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

type dbCategory struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   uint      `db:"version"`
}
