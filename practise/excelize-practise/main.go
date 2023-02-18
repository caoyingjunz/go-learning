package main

import (
	"go-learning/practise/excelize-practise/excel"
)

func main() {
	f := excel.NewExcelFile()
	defer f.Close()

	f.SetSheetName("Sheet1", "SheetName")

	_ = f.SetCellSlice([]string{"name1", "name2", "name3", "name4"})
	_ = f.SetCellSlice([]string{"1", "2", "3", "4"})

	f.SaveAs("Book1.xlsx")
}
