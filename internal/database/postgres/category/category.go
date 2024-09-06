package pgstorecategory

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"time"
)

type CategoryPostgresStore struct {
	db    *sqlx.DB
	pgxdb *pgxpool.Pool
}

func New(db *sqlx.DB, pgxdb *pgxpool.Pool) *CategoryPostgresStore {
	return &CategoryPostgresStore{
		db: db, pgxdb: pgxdb,
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

	query := "SELECT * FROM public.category"
	rows, err := db.pgxdb.Query(ctx, query)

	if rows.Err(); err != nil {
		return nil, err
	}
	categories := make([]model.Category, 0)

	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug, &category.CreatedAt, &category.UpdatedAt, &category.Version); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return categories, nil
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
