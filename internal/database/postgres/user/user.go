package pgstoreuser

import (
	"context"
	featureuser "github.com/Sanchir01/candles_backend/internal/feature/user"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgresStore struct {
	pgxdb *pgxpool.Pool
}

func New(pgxdb *pgxpool.Pool) *UserPostgresStore {
	return &UserPostgresStore{
		pgxdb: pgxdb,
	}
}

func (s *UserPostgresStore) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT id ,title,slug, phone, created_at, updated_at, version, role FROM public.users WHERE phone = $1`
	var user featureuser.DBUser
	if err := conn.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.Title, &user.Slug, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role); err != nil {
		return nil, err
	}
	return (*model.User)(&user), nil
}

func (s *UserPostgresStore) GetById(ctx context.Context, userid uuid.UUID) (*model.User, error) {
	conn, err := s.pgxdb.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT id ,title,slug, phone, created_at, updated_at, version, role FROM public.users WHERE id = $1`
	var user featureuser.DBUser
	if err := conn.QueryRow(ctx, query, userid).Scan(
		&user.ID, &user.Title, &user.Slug, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role); err != nil {
		return nil, err
	}
	return (*model.User)(&user), nil
}
