package category

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{
		primaryDB,
	}
}

func (s *Repository) CategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.category WHERE id = $1"

	var category DBCategory
	err = conn.QueryRow(ctx, query, id).Scan(&category.ID, &category.Title, category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version)
	if err != nil {
		return nil, err
	}

	return (*model.Category)(&category), nil

}
func (s *Repository) CategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.category WHERE slug = $1"

	var category DBCategory
	err = conn.QueryRow(ctx, query, slug).Scan(&category.ID, &category.Title, category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version)
	if err != nil {
		return nil, err
	}

	return (*model.Category)(&category), nil

}

func (s *Repository) AllCategories(ctx context.Context) ([]model.Category, error) {
	conn, err := s.primaryDB.Acquire(ctx)
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

func (s *Repository) CreateCategory(ctx context.Context, name, slug string) (uuid.UUID, error) {
	conn, err := s.primaryDB.Acquire(ctx)
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

func (s *Repository) UpdateCategory(ctx context.Context, id uuid.UUID, name, slug string) (uuid.UUID, error) {
	conn, err := s.primaryDB.Acquire(ctx)
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