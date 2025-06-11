package app

import (
	"github.com/Sanchir01/candles_backend/internal/feature/candles"
	"github.com/Sanchir01/candles_backend/internal/feature/category"
	"github.com/Sanchir01/candles_backend/internal/feature/color"
	"github.com/Sanchir01/candles_backend/internal/feature/events"
	"github.com/Sanchir01/candles_backend/internal/feature/order"
	"github.com/Sanchir01/candles_backend/internal/feature/user"
)

type Repositories struct {
	ColorRepository    *color.Repository
	CategoryRepository *category.Repository
	CandlesRepository  *candles.RepositoryCandles
	UserRepository     *user.RepositoryUser
	OrderRepository    *order.Repository
	EventRepository    *events.Repository
}

func NewRepositories(databases *Database) *Repositories {
	return &Repositories{
		ColorRepository:    color.NewRepository(databases.PrimaryDB),
		CategoryRepository: category.NewRepository(databases.PrimaryDB),
		CandlesRepository:  candles.NewRepositoryCandles(databases.PrimaryDB, databases.RedisDB),
		UserRepository:     user.NewRepositoryUser(databases.PrimaryDB, databases.RedisDB),
		OrderRepository:    order.NewRepository(databases.PrimaryDB),
		EventRepository:    events.NewRepository(databases.PrimaryDB),
	}
}
