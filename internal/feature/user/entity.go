package user

import (
	"github.com/google/uuid"
	"time"
)

type DBUser struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   uint      `db:"version"`
	Phone     string    `db:"phone"`
	Email     string    `json:"email"`
	Password  []byte    `db:"password"`
	Role      string    `db:"role"`
}
