package featurecategory

import (
	"github.com/google/uuid"
	"time"
)

type DBCategory struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   uint      `db:"version"`
}
