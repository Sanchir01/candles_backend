package featurecolor

import "github.com/Sanchir01/candles_backend/internal/gql/model"

func MapColorToGql(colors []model.Color) ([]*model.Color, error) {
	colorsChan := make(chan *model.Color, len(colors))
	var colorsPtrs []*model.Color
	go func() {
		for i := range colors {
			colorsChan <- &colors[i]
		}
		close(colorsChan)
	}()
	for category := range colorsChan {
		colorsPtrs = append(colorsPtrs, category)
	}
	return colorsPtrs, nil
}
