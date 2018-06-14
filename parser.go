package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

const (
	// FilePrompt : prompt for reading file path input
	FilePrompt string = "Enter xlsx file location :"
	// SheetPrompt : prompt for reading sheet input
	SheetPrompt string = "Enter the sheet name"
	// StatementError : error string for creating and checking invalid statement ids
	StatementError string = "invalid id"
)

func getNames() (string, string, error) {
	fName, err := readInputWithPrompt(FilePrompt)
	if err != nil {
		return "", "", err
	}

	sheetName, err := readInputWithPrompt(SheetPrompt)
	if err != nil {
		return "", "", err
	}
	return fName, sheetName, nil
}

func getValidCells(sheet *xlsx.Sheet, r *regexp.Regexp, index int) []*MatchedCell {
	var validCells []*MatchedCell
	for i, row := range sheet.Rows {
		if len(row.Cells) > index {
			text := strings.TrimSpace(row.Cells[index].String())
			if r.MatchString(text) {
				validCells = append(validCells, &MatchedCell{
					RowIndex:  i,
					ColIndex:  index,
					CellValue: text,
				})
			}
		}
	}

	return validCells
}

func defaultOutputCells(sheet, output *xlsx.Sheet, validCells []*MatchedCell, rowIndex, maxRows, i int) error {
	prevIndex := i - 1
	indexRow := 1
	if prevIndex > -1 {
		indexRow = validCells[prevIndex].RowIndex + 1
	}

	endRow := rowIndex + 1
	if endRow > maxRows {
		endRow = maxRows
	}

	for ; indexRow < endRow; indexRow++ {
		if err := copyCells(sheet, output, validCells, indexRow); err != nil {
			return err
		}
	}
	return nil
}

func monthOutputCells(sheet, output *xlsx.Sheet, validCells []*MatchedCell, sz, rowIndex, i int) error {
	endRow := sheet.MaxRow
	nextIndex := i + 1
	if sz > nextIndex {
		endRow = validCells[nextIndex].RowIndex
	}

	for j := rowIndex; j < endRow; j++ {
		outputRow := output.AddRow()
		for _, cell := range sheet.Rows[j].Cells {
			outputCell := outputRow.AddCell()
			if err := copyCellWithFill(cell, outputCell, MonthReportFGColour); err != nil {
				return err
			}
		}
	}
	return nil
}

func statementOutputCells(sheet, output *xlsx.Sheet, validCells []*MatchedCell, sz, rowIndex, i int) error {
	endRow := sheet.MaxRow
	nextIndex := i + 1
	if nextIndex < sz {
		endRow = validCells[nextIndex].RowIndex - 1
	}

	for j := rowIndex; j < endRow; j++ {
		if err := copyCells(sheet, output, validCells, j); err != nil {
			return err
		}
	}
	return nil
}

func creditOutputCells(sheet *xlsx.Sheet, validCells []*MatchedCell, sz int, directoryName string) error {
	for i, j := 0, 0; i < sz; {
		file := xlsx.NewFile()
		output, err := file.AddSheet(validCells[i].CellValue)
		if err != nil {
			return err
		}
		if err := copyCells(sheet, output, validCells, 0); err != nil {
			return err
		}
		for j = i; j < sz; j++ {
			if validCells[i].CellValue != validCells[j].CellValue {
				break
			}
			if err := copyCells(sheet, output, validCells, validCells[j].RowIndex); err != nil {
				return err
			}
		}
		file.Save(getOutputName(directoryName, validCells[i].CellValue, Ext))
		i = j
	}
	return nil
}

func outputCells(sheet, output *xlsx.Sheet, out string, cell *MatchedCell, validCells []*MatchedCell, i int) error {
	if err := copyCells(sheet, output, validCells, 0); err != nil {
		return err
	}

	if out == Annual {
		if err := defaultOutputCells(sheet, output, validCells, cell.RowIndex, sheet.MaxRow, i); err != nil {
			return err
		}
	} else if out == Month {
		if err := monthOutputCells(sheet, output, validCells, len(validCells), cell.RowIndex, i); err != nil {
			return err
		}
	} else if out == Shift {
		if err := copyCells(sheet, output, validCells, 1); err != nil {
			return err
		}
		if err := defaultOutputCells(sheet, output, validCells, cell.RowIndex, sheet.MaxRow, i); err != nil {
			return err
		}
	} else if out == Statement {
		if err := statementOutputCells(sheet, output, validCells, len(validCells), cell.RowIndex, i); err != nil {
			return err
		}
	}
	return nil
}

func saveCredit(sheet *xlsx.Sheet, validCells []*MatchedCell, directoryName string) error {
	sort.Slice(validCells, func(i, j int) bool {
		// Regex is used before adding, therefore error check can be avoided
		a, _ := strconv.Atoi(validCells[i].CellValue)
		b, _ := strconv.Atoi(validCells[j].CellValue)

		return a < b
	})

	if err := creditOutputCells(sheet, validCells, len(validCells), directoryName); err != nil {
		return err
	}
	return nil
}

func saveData(sheet *xlsx.Sheet, validCells []*MatchedCell, directoryName, out string) error {
	if out == CreditCards {
		err := saveCredit(sheet, validCells, directoryName)
		if err != nil {
			return err
		}
		return nil
	}

	for i, cell := range validCells {
		file := xlsx.NewFile()
		output, err := file.AddSheet(cell.CellValue)
		if err != nil {
			return err
		}

		if err := outputCells(sheet, output, out, cell, validCells, i); err != nil {
			if err.Error() == StatementError {
				continue
			}
			return err
		}

		file.Save(getOutputName(directoryName, cell.CellValue, Ext))
	}
	return nil
}

func parse(fName, basename, sheetName, regex, out string, colName rune) error {
	if fName == basename {
		fName = fmt.Sprintf("%s.%s", fName, Ext)
	}

	xlFile, err := xlsx.OpenFile(fName)
	if err != nil {
		fmt.Printf("unable to file :%s\nperhaps it does not exist\n", fName)
		return err
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}

	sheet, err := getSheet(xlFile, sheetName)
	if err != nil {
		fmt.Printf("unable to open sheet :%s\nperhaps it does not exist\n", sheetName)
		return err
	}

	index, err := convertRuneToIndex(colName)
	if err != nil {
		return err
	}

	validCells := getValidCells(sheet, r, index)
	sz := len(validCells)
	directoryName := fmt.Sprintf("%s_output", basename)
	if sz > 0 {
		if err := createDirIfNotExists(directoryName); err != nil {
			return err
		}
	}

	if err = saveData(sheet, validCells, directoryName, out); err != nil {
		return err
	}
	fmt.Printf("successfully parsed %s\ncheck out %s for files", fName, directoryName)
	return nil
}
