package user

import (
	"context"
	"encoding/base64"
	"fmt"
	sq "github.com/Masterminds/squirrel"
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
	passwordBase64 := base64.StdEncoding.EncodeToString(user.Password)
	return &model.User{
		ID:        user.ID,
		Title:     user.Title,
		Slug:      user.Slug,
		Phone:     user.Phone,
		Role:      user.Role,
		Email:     user.Email,
		Version:   user.Version,
		Password:  passwordBase64,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, arg, err := sq.
		Select("id ,title,slug, phone, created_at, updated_at, version, role, password").
		From("public.users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var user DBUser
	if err := conn.QueryRow(ctx, query, arg...).Scan(
		&user.ID, &user.Title, &user.Slug, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role, &user.Password); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("Неправильный логин или пароль")
		}
		return nil, err
	}

	return &model.User{
		ID:        user.ID,
		Title:     user.Title,
		Slug:      user.Slug,
		Phone:     user.Phone,
		Role:      user.Role,
		Email:     user.Email,
		Version:   user.Version,
		Password:  "",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
func (r *Repository) GetBySlug(ctx context.Context, slug string) (*model.User, error) {
	conn, err := r.primartDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, arg, err := sq.
		Select("id ,title,slug, phone, created_at, updated_at, version, role, password").
		From("public.users").
		Where(sq.Eq{"slug": slug}).
		PlaceholderFormat(sq.Dollar).ToSql()
	var user DBUser

	if err := conn.QueryRow(ctx, query, arg...).Scan(
		&user.ID, &user.Title, &user.Slug, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role, &user.Password); err != nil {
		return nil, err
	}

	return &model.User{
		ID:        user.ID,
		Title:     user.Title,
		Slug:      user.Slug,
		Phone:     user.Phone,
		Role:      user.Role,
		Email:     user.Email,
		Version:   user.Version,
		Password:  "",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
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

	return &model.User{
		ID:        user.ID,
		Title:     user.Title,
		Slug:      user.Slug,
		Phone:     user.Phone,
		Role:      user.Role,
		Email:     user.Email,
		Version:   user.Version,
		Password:  "",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *Repository) CreateUser(ctx context.Context, title, phone, slug, email, role string, password []byte, tx pgx.Tx) (*model.User, error) {
	query, arg, err := sq.
		Insert("users").
		Columns("title", "slug", "phone", "email", "role", "password").
		Values(title, slug, phone, email, role, password).
		Suffix("RETURNING id, phone, role, email").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}
	var users DBUser

	if err := tx.QueryRow(ctx, query, arg...).Scan(&users.ID, &users.Phone, &users.Role, &users.Email); err != nil {
		if err == pgx.ErrTxCommitRollback {
			return nil, fmt.Errorf("ошибка при создании пользователя")
		}
		return nil, err
	}
	passwordBase64 := base64.StdEncoding.EncodeToString(users.Password)
	return &model.User{
		ID:        users.ID,
		Title:     users.Title,
		Slug:      users.Slug,
		Phone:     users.Phone,
		Role:      users.Role,
		Email:     users.Email,
		Version:   users.Version,
		Password:  passwordBase64,
		CreatedAt: users.CreatedAt,
		UpdatedAt: users.UpdatedAt,
	}, nil
}
