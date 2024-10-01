package color

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primartDB *pgxpool.Pool
}

func NewRepository(primartDB *pgxpool.Pool) *Repository {
	return &Repository{primartDB: primartDB}
}

func (r *Repository) AllColor(ctx context.Context) ([]model.Color, error) {
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query := "SELECT id , title, slug, created_at, updated_at, version FROM public.color"
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
	query := `INSERT INTO color(title,slug) VALUES ($1,$2) RETURNING id`

	var id uuid.UUID

	if err = conn.QueryRow(ctx, query, title, slug).Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
