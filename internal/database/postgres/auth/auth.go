package pgstoreauth

import (
	"context"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type AuthPostgresStore struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *AuthPostgresStore {
	return &AuthPostgresStore{
		db: db,
	}
}

func (s *AuthPostgresStore) Register(ctx context.Context, title, slug, phone, role string) (*model.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query := `INSERT INTO users (title, slug, phone, role) 
          VALUES ($1, $2, $3, $4) 
          RETURNING id, phone, role`

	var users dbUser

	if err := conn.QueryRow(ctx, query, title, slug, phone, role).Scan(&users.ID, &users.Phone, &users.Role); err != nil {
		return nil, err
	}

	return (*model.User)(&users), nil
}

func (s *AuthPostgresStore) Login(ctx context.Context, phone string) (*model.User, error) {
	conn, err := s.db.Acquire(ctx)
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
