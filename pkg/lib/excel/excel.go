package excel

import (
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	"github.com/xuri/excelize/v2"
)

func CreateFileExcel(input model.CreateOrderInput) error {
	f := excelize.NewFile()

	sheetName := "Order"
	f.SetCellValue(sheetName, "B1", "Название товара")
	f.SetCellValue(sheetName, "C1", "Количество")
	f.SetCellValue(sheetName, "D1", "Цена")
	row := 2
	// Set active sheet of the workbook.
	for _, item := range input.Items {
		f.SetCellValue(sheetName, "B"+string(row), item.ProductsID)
		f.SetCellValue(sheetName, "C"+string(row), item.Quantity)
		f.SetCellValue(sheetName, "D"+string(row), item.Price)
		row++
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Order.xlsx"); err != nil {
		return err
	}
	return nil
}
