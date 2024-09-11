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
	query := `INSERT INTO users(title,slug, phone,role) VALUES ($1,$2,$3,$4) RETURNING id phone role`

	var users dbUser

	if err := conn.QueryRow(ctx, query, title, slug, phone, role).Scan(&users.ID, &users.Phone, &users.Role); err != nil {
		return nil, err
	}
	return (*model.User)(&users), nil
}


func (s *AuthPostgresStore) Login(ctx context.Context) (*model.User, error) {
  conn, err := s.db.Acquire(ctx) 
  if err !=nil {
    return nil, error
  } 
  defer c
}
type dbUser struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	Slug      string    `db:"slug"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   uint      `db:"version"` Role      string    `db:"role"`
}
