package category

import "github.com/Sanchir01/candles_backend/internal/gql/model"

func MapCategoryToGql(c []model.Category) (categories []*model.Category, err error) {
	categoriesChan := make(chan *model.Category, len(c))
	var categoryPtrs []*model.Category
	go func() {
		for i := range c {
			categoriesChan <- &c[i]
		}
		close(categoriesChan)
	}()

	for category := range categoriesChan {
		categoryPtrs = append(categoryPtrs, category)
	}
	return categoryPtrs, nil
}
