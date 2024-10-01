package pgstoreuser

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
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
	var user dbUser
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
	var user dbUser
	if err := conn.QueryRow(ctx, query, userid).Scan(
		&user.ID, &user.Title, &user.Slug, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role); err != nil {
		return nil, err
	}
	return (*model.User)(&user), nil
}

type dbUser struct {
	ID        uuid.UUID  `db:"id"`
	Title     string     `db:"title"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	Slug      string     `db:"slug"`
	Version   uint       `db:"version"`
	Phone     string     `db:"phone"`
	Role      model.Role `db:"role"`
}
