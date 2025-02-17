package app

import (
	"github.com/Sanchir01/candles_backend/internal/feature/candles"
	"github.com/Sanchir01/candles_backend/internal/feature/category"
	"github.com/Sanchir01/candles_backend/internal/feature/color"
	"github.com/Sanchir01/candles_backend/internal/feature/order"
	"github.com/Sanchir01/candles_backend/internal/feature/user"
)

type Services struct {
	ColorService    *color.Service
	CategoryService *category.Service
	CandlesService  *candles.ServiceCandles
	UserService     *user.Service
	OrderService    *order.Service
}

func NewServices(repos *Repositories, storages *Storages) *Services {
	return &Services{
		ColorService:    color.NewService(repos.ColorRepository),
		CategoryService: category.NewService(repos.CategoryRepository),
		CandlesService:  candles.NewServiceCandles(repos.CandlesRepository, storages.CandlesStorage),
		UserService:     user.NewService(repos.UserRepository),
		OrderService:    order.NewService(repos.OrderRepository),
	}
}
