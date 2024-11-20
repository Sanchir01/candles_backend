package color

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Repository struct {
	primartDB *pgxpool.Pool
}

func NewRepository(primartDB *pgxpool.Pool) *Repository {
	return &Repository{
		primartDB,
	}
}

func (r *Repository) AllColor(ctx context.Context) ([]model.Color, error) {
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, _, err :=
		sq.
			Select("id , title, slug, created_at, updated_at,version").
			From("public.color").
			ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, query)
	defer rows.Close()

	colors := make([]model.Color, 0)

	for rows.Next() {
		var color model.Color
		if err := rows.Scan(&color.ID, &color.Title, &color.Slug, &color.CreatedAt, &color.UpdatedAt, &color.Version); err != nil {
			return nil, err
		}
		colors = append(colors, color)
	}
	return colors, nil
}

func (r *Repository) CreateColor(ctx context.Context, title, slug string) (uuid.UUID, error) {
	conn, err := r.primartDB.Acquire(ctx)
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
	if err = conn.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *Repository) ColorByManyId(ctx context.Context, ids []uuid.UUID) ([]*model.Color, error) {
	log.Printf("DataLoader keys many load sanchir  test: %v", ids)
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := "SELECT id, title, slug, version, created_at, updated_at FROM public.color WHERE id = ANY($1)"

	rows, err := conn.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	var colors []*model.Color
	for rows.Next() {
		var color model.Color
		if err := rows.Scan(&color.ID, &color.Title, &color.Slug, &color.Version, &color.CreatedAt, &color.UpdatedAt); err != nil {
			return nil, err
		}
		colors = append(colors, &color)
	}
	return colors, nil
}

func (r *Repository) ColorById(ctx context.Context, id uuid.UUID) (*model.Color, error) {
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, _, err := sq.Select("id , title, slug, created_at, updated_at,version ").From("public.color").
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	var color DBColor

	if err := conn.QueryRow(ctx, query, id).Scan(
		&color.ID, &color.Title, &color.Slug, &color.Version, &color.CreatedAt, &color.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return (*model.Color)(&color), nil

}
func (r *Repository) ColorBySlug(ctx context.Context, slug string) (*model.Color, error) {
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, _, err := sq.Select("id , title, slug, created_at, updated_at,version").From("public.color").
		Where(sq.Eq{"slug": slug}).PlaceholderFormat(sq.Dollar).ToSql()
	var color DBColor

	if err := conn.QueryRow(ctx, query, slug).Scan(
		&color.ID, &color.Title, &color.Slug, &color.Version, &color.CreatedAt, &color.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return (*model.Color)(&color), nil
}
