package pgstorecategory

import (
	"context"
	"time"

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
func (s *CategoryPostgresStore) CategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.category WHERE id = $1"

	var category dbCategory
	err = conn.QueryRow(ctx, query, id).Scan(&category.ID, &category.Title, category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version)
	if err != nil {
		return nil, err
	}

	return (*model.Category)(&category), nil

}
func (s *CategoryPostgresStore) CategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.category WHERE slug = $1"

	var category dbCategory
	err = conn.QueryRow(ctx, query, slug).Scan(&category.ID, &category.Title, category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version)
	if err != nil {
		return nil, err
	}

	return (*model.Category)(&category), nil

}

func (s *CategoryPostgresStore) AllCategories(ctx context.Context) ([]model.Category, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.category"
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
	query := "INSERT INTO category(title, slug) VALUES($1, $2) RETURNING id"
	var id uuid.UUID

	row := conn.QueryRow(ctx, query, name, slug)

	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *CategoryPostgresStore) UpdateCategory(ctx context.Context, id uuid.UUID, name, slug string) (uuid.UUID, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()
	var idReturning uuid.UUID
	query := "PDATE category SET title = $1, slug = $2 WHERE id = $3"
	row := conn.QueryRow(ctx, query, name, slug)
	if err := row.Scan(&idReturning); err != nil {
		return uuid.Nil, err
	}
	return idReturning, nil
}

type dbCategory struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   uint      `db:"version"`
}
