package app

import (
	featurecolor "github.com/Sanchir01/candles_backend/internal/feature/color"
)

const (
	colorLoaderByIDMaxBatch               int = 100
	articleBlockLoaderByArticleIDMaxBatch int = 10
)

type DataLoaders struct {
	ColorDataLoader *featurecolor.ColorDataLoader
}

func NewDataLoaders(repos *Repositories) *DataLoaders {
	return &DataLoaders{
		ColorDataLoader: featurecolor.NewDataLoader(repos.ColorRepository, colorLoaderByIDMaxBatch),
	}
}
