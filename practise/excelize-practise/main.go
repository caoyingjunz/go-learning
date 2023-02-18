package main

import (
	"fmt"
	"sync"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := NewExcelFile()
	defer f.Close()

	f.SetSheetName("Sheet1", "SheetName")

	_ = f.SetCellSlice([]string{"name1", "name2", "name3", "name4"})
	_ = f.SetCellSlice([]string{"1", "2", "3"})

	f.SaveAs("Book1.xlsx")
}

type ExcelFile struct {
	lock sync.RWMutex

	sheetName string
	f         *excelize.File
	cursor    int // 游标，用户标记当前所处位置
}

func NewExcelFile() *ExcelFile {
	return &ExcelFile{
		f:      excelize.NewFile(),
		cursor: 1,
	}
}

func (file *ExcelFile) SetSheetName(oldName, NewName string) {
	file.sheetName = NewName
	_ = file.f.SetSheetName(oldName, NewName)
}

// SetCellSlice 以行为单位
func (file *ExcelFile) SetCellSlice(cells []string) error {
	file.lock.Lock()
	defer file.lock.Unlock()

	if len(cells) == 0 {
		return nil
	}
	for index, cell := range cells {
		excelIndex := fmt.Sprintf("%s%d", file.parseExcelIndex(index), file.cursor)
		if err := file.f.SetCellValue(file.sheetName, excelIndex, cell); err != nil {
			return err
		}
	}

	file.cursor += 1
	return nil
}

func (file *ExcelFile) parseExcelIndex(i int) string {
	index := 65 + i
	indexRune := rune(index)

	return string(indexRune)
}

func (file *ExcelFile) SaveAs(f string) error {
	return file.f.SaveAs(f)
}

func (file *ExcelFile) Close() error {
	return file.f.Close()
}
