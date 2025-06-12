package app

import (
	"github.com/Sanchir01/candles_backend/internal/feature/candles"
	"github.com/Sanchir01/candles_backend/internal/feature/category"
	"github.com/Sanchir01/candles_backend/internal/feature/color"
	"github.com/Sanchir01/candles_backend/internal/feature/events"
	"github.com/Sanchir01/candles_backend/internal/feature/order"
	"github.com/Sanchir01/candles_backend/internal/feature/user"
	"log/slog"
)

type Services struct {
	ColorService    *color.Service
	CategoryService *category.Service
	CandlesService  *candles.ServiceCandles
	UserService     *user.Service
	OrderService    *order.Service
	EventService    *events.EventService
}

func NewServices(repos *Repositories, storages *Storages, db *Database, kaf *Producer, l *slog.Logger) *Services {
	return &Services{
		ColorService:    color.NewService(repos.ColorRepository, repos.EventRepository, db.PrimaryDB),
		CategoryService: category.NewService(repos.CategoryRepository),
		CandlesService:  candles.NewServiceCandles(repos.CandlesRepository, storages.CandlesStorage),
		UserService:     user.NewService(repos.UserRepository),
		OrderService:    order.NewService(repos.OrderRepository, repos.EventRepository),
		EventService:    events.NewEventService(l, repos.EventRepository, kaf),
	}
}
