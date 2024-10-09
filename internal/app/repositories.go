package app

import (
	"github.com/Sanchir01/candles_backend/internal/feature/candles"
	"github.com/Sanchir01/candles_backend/internal/feature/category"
	"github.com/Sanchir01/candles_backend/internal/feature/color"
	"github.com/Sanchir01/candles_backend/internal/feature/order"
	"github.com/Sanchir01/candles_backend/internal/feature/user"
)

type Repositories struct {
	ColorRepository    *color.Repository
	CategoryRepository *category.Repository
	CandlesRepository  *candles.Repository
	UserRepository     *user.Repository
	OrderRepository    *order.Repository
}

func NewRepositories(databases *Database) *Repositories {
	return &Repositories{
		ColorRepository:    color.NewRepository(databases.PrimaryDB),
		CategoryRepository: category.NewRepository(databases.PrimaryDB),
		CandlesRepository:  candles.NewRepository(databases.PrimaryDB),
		UserRepository:     user.NewRepository(databases.PrimaryDB),
		OrderRepository:    order.NewRepository(databases.PrimaryDB),
	}
}
