package pgstoreauth

import (
	"context"
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type AuthStorePosrtgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *AuthStorePosrtgres {
	return &AuthStorePosrtgres{
		db: db,
	}
}

func Register(ctx context.Context, title, phone, slug, role string, tx pgx.Tx) (*model.User, error) {
	query := `INSERT INTO users (title, slug, phone, role)
	     VALUES ($1, $2, $3, $4)
	     RETURNING id, phone, role`

	var users dbUser

	if err := tx.QueryRow(ctx, query, title, slug, phone, role).Scan(&users.ID, &users.Phone, &users.Role); err != nil {
		if err == pgx.ErrTxCommitRollback {
			return nil, fmt.Errorf("ошибка при создании пользователя")
		}
		return nil, err
	}

	return (*model.User)(&users), nil
}

func (s *AuthStorePosrtgres) Login(ctx context.Context, phone string) (*model.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query := `SELECT id ,title,slug, phone, created_at, updated_at, version, role FROM public.users WHERE phone = $1`
	var user dbUser

	if err := conn.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.Title, &user.Slug, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("пользователь с номером телефона %s не найден", phone)
		}
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
