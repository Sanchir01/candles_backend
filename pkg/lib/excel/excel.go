package excel

import (
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/xuri/excelize/v2"
)

func CreateFileExcel(candles []model.Candles, quantity []int, filePath, sheetName string) error {
	if len(candles) != len(quantity) {
		return fmt.Errorf("length of candles and quantity do not match")
	}
	f := excelize.NewFile()

	f.SetSheetName(f.GetSheetName(0), sheetName)
	headers := []string{"Название", "Цена", "Количесво", "Цена за эти товары"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+i))
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return fmt.Errorf("failed to set header: %v", err)
		}
	}
	maxLength := make([]int, len(headers))

	for rowIndex, product := range candles {
		row := rowIndex + 2
		err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.Slug)
		if err != nil {
			return fmt.Errorf("failed to set title: %v", err)
		}
		maxLength[0] = max(maxLength[0], len(product.Slug))
		err = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.Price)
		if err != nil {
			return fmt.Errorf("failed to set price: %v", err)
		}
		maxLength[1] = max(maxLength[1], len(fmt.Sprintf("%f", product.Price)))

		err = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), quantity[rowIndex])
		if err != nil {
			return fmt.Errorf("failed to set quantity: %v", err)
		}

		maxLength[2] = max(maxLength[2], len(fmt.Sprintf("%d", product.Slug)))

		totalPrice := quantity[rowIndex] * product.Price
		err = f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), totalPrice)
		if err != nil {
			return fmt.Errorf("failed to set product quantity price: %v", err)
		}

		maxLength[3] = max(maxLength[3], len(fmt.Sprintf("%f", totalPrice)))
	}
	for i, maxLength := range maxLength {
		col := string('A' + i)
		width := float64(maxLength + 2)
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return fmt.Errorf("failed to set column width for %s: %v", col, err)
		}
	}
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save excel file: %v", err)
	}

	return nil
}
