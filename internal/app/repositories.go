package app

import (
	"github.com/Sanchir01/candles_backend/internal/feature/candles"
	"github.com/Sanchir01/candles_backend/internal/feature/category"
	"github.com/Sanchir01/candles_backend/internal/feature/color"
)

type Repositories struct {
	ColorRepository    *color.Repository
	CategoryRepository *category.Repository
	CandlesReository   *candles.Repository
}

func NewRepositories(databases *Database) *Repositories {
	return &Repositories{
		ColorRepository:    color.NewRepository(databases.PrimaryDB),
		CategoryRepository: category.NewRepository(databases.PrimaryDB),
		CandlesReository:   candles.NewRepository(databases.PrimaryDB),
	}
}
