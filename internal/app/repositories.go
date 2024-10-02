package app

import (
	"github.com/Sanchir01/candles_backend/internal/feature/category"
	"github.com/Sanchir01/candles_backend/internal/feature/color"
)

type Repositories struct {
	ColorRepository    *color.Repository
	CategoryRepository *category.Repository
}

func newRepositories(databases *Database) *Repositories {
	return &Repositories{
		ColorRepository:    color.NewRepository(databases.PrimaryDB),
		CategoryRepository: category.NewRepository(databases.PrimaryDB),
	}
}
