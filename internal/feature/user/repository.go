package user

import (
	"context"
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primartDB *pgxpool.Pool
}

func NewRepository(primartDB *pgxpool.Pool) *Repository {
	return &Repository{
		primartDB,
	}
}
func (r *Repository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT id ,title,slug, phone, created_at, updated_at, version, role FROM public.users WHERE phone = $1`
	var user DBUser
	if err := conn.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.Title, &user.Slug, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role); err != nil {
		return nil, err
	}
	return (*model.User)(&user), nil
}

func (r *Repository) GetById(ctx context.Context, userid uuid.UUID) (*model.User, error) {
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT id ,title,slug, phone, created_at, updated_at, version, role FROM public.users WHERE id = $1`
	var user DBUser
	if err := conn.QueryRow(ctx, query, userid).Scan(
		&user.ID, &user.Title, &user.Slug, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role); err != nil {
		return nil, err
	}
	return (*model.User)(&user), nil
}

func (r *Repository) CreateUser(ctx context.Context, title, phone, slug, role string, tx pgx.Tx) (*model.User, error) {
	query := `INSERT INTO users (title, slug, phone, role)
	     VALUES ($1, $2, $3, $4)
	     RETURNING id, phone, role`

	var users DBUser

	if err := tx.QueryRow(ctx, query, title, slug, phone, role).Scan(&users.ID, &users.Phone, &users.Role); err != nil {
		if err == pgx.ErrTxCommitRollback {
			return nil, fmt.Errorf("ошибка при создании пользователя")
		}
		return nil, err
	}

	return (*model.User)(&users), nil
}
