package candles

import "github.com/Sanchir01/candles_backend/internal/gql/model"

func MapCandlesToGql(candles []model.Candles) ([]*model.Candles, error) {
	candlesChan := make(chan *model.Candles, len(candles))
	var candlesPtrs []*model.Candles
	go func() {
		for i := range candles {
			candlesChan <- &candles[i]
		}
		close(candlesChan)
	}()

	for candlesgql := range candlesChan {
		candlesPtrs = append(candlesPtrs, candlesgql)
	}
	return candlesPtrs, nil
}
