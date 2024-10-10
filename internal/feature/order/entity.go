package order

import (
	"github.com/google/uuid"
	"time"
)

type DBOrderItems struct {
}

type DBOrders struct {
	ID          uuid.UUID `db:"id"`
	CreatedAt   time.Time `db:"createdAt"`
	UpdatedAt   time.Time `db:"updatedAt"`
	Status      string    `db:"status"`
	UserID      uuid.UUID `db:"userId"`
	TotalAmount int       `db:"total_amount"`
	Version     uint      `db:"version"`
}

type OrderItemsInput struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}
