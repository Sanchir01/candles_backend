package candles

import (
	"github.com/google/uuid"
	"time"
)

type DBCandles struct {
	ID         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	Slug       string    `db:"slug"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Version    uint      `db:"version"`
	Price      int       `db:"price"`
	Images     []string  `db:"images"`
	ColorID    uuid.UUID `db:"color_id"`
	CategoryID uuid.UUID `db:"category_id"`
}
