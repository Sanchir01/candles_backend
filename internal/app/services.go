package app

import (
	"github.com/Sanchir01/candles_backend/internal/bot"
	"github.com/Sanchir01/candles_backend/internal/feature/candles"
	"github.com/Sanchir01/candles_backend/internal/feature/category"
	"github.com/Sanchir01/candles_backend/internal/feature/color"
	"github.com/Sanchir01/candles_backend/internal/feature/order"
	"github.com/Sanchir01/candles_backend/internal/feature/user"
)

type Services struct {
	ColorService    *color.Service
	CategoryService *category.Service
	CandlesService  *candles.Service
	UserService     *user.Service
	OrderService    *order.Service
}

func NewServices(repos *Repositories, storages *Storages, bot *bot.Bot) *Services {
	return &Services{
		ColorService:    color.NewService(repos.ColorRepository),
		CategoryService: category.NewService(repos.CategoryRepository),
		CandlesService:  candles.NewService(repos.CandlesRepository, storages.CandlesStorage),
		UserService:     user.NewService(repos.UserRepository),
		OrderService:    order.NewService(repos.OrderRepository, bot),
	}
}
