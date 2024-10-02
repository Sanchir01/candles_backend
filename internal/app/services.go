package app

import (
	"github.com/Sanchir01/candles_backend/internal/feature/category"
	"github.com/Sanchir01/candles_backend/internal/feature/color"
)

type Services struct {
	ColorService    *color.Service
	CategoryService *category.Service
}

func NewServices(repos *Repositories, storages *Storages) *Services {
	return &Services{
		ColorService:    color.NewService(repos.ColorRepository),
		CategoryService: category.NewService(repos.CategoryRepository),
	}
}
