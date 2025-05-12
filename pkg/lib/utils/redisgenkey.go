package utils

import (
	"fmt"
	"strings"

	"github.com/Sanchir01/candles_backend/internal/gql/model"
)

func GenerateCacheKey(sort *model.CandlesSortEnum, filter *model.CandlesFilterInput, pageSize uint, pageNumber uint) string {
	keyParts := []string{"candles"}

	if sort != nil {
		keyParts = append(keyParts, fmt.Sprintf("sort:%v", *sort))
	}

	if filter != nil {
		if filter.CategoryID != nil {
			keyParts = append(keyParts, fmt.Sprintf("cat:%d", *filter.CategoryID))
		}
		if filter.ColorID != nil {
			keyParts = append(keyParts, fmt.Sprintf("col:%d", *filter.ColorID))
		}
	}

	keyParts = append(keyParts, fmt.Sprintf("pageSize:%d", pageSize))
	keyParts = append(keyParts, fmt.Sprintf("page:%d", pageNumber))

	return strings.Join(keyParts, "|")
}
