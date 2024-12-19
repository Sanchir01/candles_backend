package excel

import (
	"fmt"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/xuri/excelize/v2"
)

func CreateFileExcel(candles []model.Candles, quantity []int, filePath string) error {
	if len(candles) != len(quantity) {
		return fmt.Errorf("length of candles and quantity do not match")
	}
	f := excelize.NewFile()

	sheetName := "Order"
	f.SetSheetName(f.GetSheetName(0), sheetName)
	headers := []string{"Title", "Price", "Quantity"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+i)) // A1, B1, C1
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return fmt.Errorf("failed to set header: %v", err)
		}
	}

	// Set active sheet of the workbook.
	for rowIndex, product := range candles {
		row := rowIndex + 2                                                     // Данные начинаются со второй строки
		err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.Slug) // Название продукта
		if err != nil {
			return fmt.Errorf("failed to set title: %v", err)
		}

		err = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.Price) // Цена
		if err != nil {
			return fmt.Errorf("failed to set price: %v", err)
		}

		err = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), quantity[rowIndex]) // Количество
		if err != nil {
			return fmt.Errorf("failed to set quantity: %v", err)
		}
	}

	// Сохраняем файл
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save excel file: %v", err)
	}

	return nil
}
