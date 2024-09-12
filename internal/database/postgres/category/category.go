package pgstorecategory

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryPostgresStore struct {
	pgxdb *pgxpool.Pool
}

func New(pgxdb *pgxpool.Pool) *CategoryPostgresStore {
	return &CategoryPostgresStore{
		pgxdb: pgxdb,
	}
}
func (s *CategoryPostgresStore) CategoryBySlug(ctx context.Context, slug string) (model.Category, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return model.Category{}, err
	}
	defer conn.Release()

	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.category.sql WHERE slug = $1"

	var category model.Category
	err = conn.QueryRow(ctx, query, slug).Scan(&category.ID, &category.Title, category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version)
	if err != nil {
		return model.Category{}, err
	}

	return category, nil

}

func (s *CategoryPostgresStore) AllCategories(ctx context.Context) ([]model.Category, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.category.sql"
	rows, err := conn.Query(ctx, query)

	if rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)

	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Title, &category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryPostgresStore) CreateCategory(ctx context.Context, name, slug string) (uuid.UUID, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()
	query := "INSERT INTO category.sql(title, slug) VALUES($1, $2) RETURNING id"
	var id uuid.UUID

	row := conn.QueryRow(ctx, query, name, slug)

	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
