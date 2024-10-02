package category

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	primartDB *pgxpool.Pool
}

func NewRepository(primartDB *pgxpool.Pool) *Repository {
	return &Repository{
		primartDB,
	}
}
