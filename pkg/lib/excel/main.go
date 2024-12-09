package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func CreateFileExcel() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	index, err := f.NewSheet("Sheet2")
	if err != nil {
		return err
	}
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		return err
	}
	return nil
}
