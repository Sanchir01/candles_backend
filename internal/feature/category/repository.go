package category

import (
	"context"
	sq "github.com/Masterminds/squirrel"
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
	query, args, err := sq.Select("id", "title", "slug", "created_at", "updated_at", "version").
		From("public.category").
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).
		ToSql()
	var category DBCategory
	err = conn.QueryRow(ctx, query, args...).Scan(&category.ID, &category.Title, category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version)
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

	query, args, err := sq.Select("id", "title", "slug", "created_at", "updated_at", "version").
		From("public.category").
		Where(sq.Eq{"slug": slug}).PlaceholderFormat(sq.Dollar).
		ToSql()

	var category DBCategory
	err = conn.QueryRow(ctx, query, args...).Scan(&category.ID, &category.Title, category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version)
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

	query, _, err := sq.Select("id , title, slug, created_at, updated_at, version ").From("public.category").ToSql()
	if err != nil {
		return nil, err
	}

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

func (s *Repository) CreateCategory(ctx context.Context, title, slug string) (uuid.UUID, error) {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()
	query, args, err := sq.Insert("color").
		Columns("title", "slug").
		Values(title, slug).
		Suffix("RETURNING id").PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID

	row := conn.QueryRow(ctx, query, args...)

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
