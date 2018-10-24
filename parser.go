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
	// filePrompt : prompt for reading file path input
	filePrompt string = "Enter xlsx file location :"
	// sheetPrompt : prompt for reading sheet input
	sheetPrompt string = "Enter the sheet name"
	// statementError : error string for creating and checking invalid statement ids
	statementError string = "invalid id"
)

func getNames() (string, string, error) {
	fName, err := readInputWithPrompt(filePrompt)
	if err != nil {
		return "", "", err
	}

	sheetName, err := readInputWithPrompt(sheetPrompt)
	if err != nil {
		return "", "", err
	}
	return fName, sheetName, nil
}

func getValidCells(sheet *xlsx.Sheet, r *regexp.Regexp, index int) []*matchedCell {
	var validCells []*matchedCell
	for i, row := range sheet.Rows {
		if len(row.Cells) > index {
			text := strings.TrimSpace(row.Cells[index].String())
			if r.MatchString(text) {
				validCells = append(validCells, &matchedCell{
					RowIndex:  i,
					ColIndex:  index,
					CellValue: text,
				})
			}
		}
	}

	return validCells
}

func defaultOutputCells(sheet, output *xlsx.Sheet, validCells []*matchedCell, rowIndex, maxRows, i int) error {
	prevIndex := i - 1
	indexRow := 1
	if prevIndex > -1 {
		indexRow = validCells[prevIndex].RowIndex + 1
	}

	endRow := rowIndex + 1
	if endRow > maxRows {
		endRow = maxRows
	}

	for ;indexRow < endRow; indexRow++ {
		if err := copyCells(sheet, output, validCells, indexRow); err != nil {
			return err
		}
	}
	return nil
}

func monthOutputCells(sheet, output *xlsx.Sheet, validCells []*matchedCell, sz, rowIndex, i int) error {
	endRow := sheet.MaxRow
	nextIndex := i + 1
	if sz > nextIndex {
		endRow = validCells[nextIndex].RowIndex
	}

	for j := rowIndex; j < endRow; j++ {
		outputRow := output.AddRow()
		for _, cell := range sheet.Rows[j].Cells {
			outputCell := outputRow.AddCell()
			if err := copyCellWithFill(cell, outputCell, monthReportFGColour); err != nil {
				return err
			}
		}
	}
	return nil
}

func statementOutputCells(sheet, output *xlsx.Sheet, validCells []*matchedCell, sz, rowIndex, i int) error {
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

func creditOutputCells(sheet *xlsx.Sheet, validCells []*matchedCell, sz int, directoryName string) error {
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
		file.Save(getOutputName(directoryName, validCells[i].CellValue, ext))
		i = j
	}
	return nil
}

func outputCells(sheet, output *xlsx.Sheet, out string, cell *matchedCell, validCells []*matchedCell, i int) error {
	// TODO set row height
	if err := copyCells(sheet, output, validCells, 0); err != nil {
		return err
	}

	switch out {
		case annual, shift:
			return defaultOutputCells(sheet, output, validCells, cell.RowIndex, sheet.MaxRow, i)
		case month:
			return monthOutputCells(sheet, output, validCells, len(validCells), cell.RowIndex, i)
		case statement:
			return statementOutputCells(sheet, output, validCells, len(validCells), cell.RowIndex, i)
		default:
			return fmt.Errorf("%s is not a valid document type", out)
	}
}

func saveCredit(sheet *xlsx.Sheet, validCells []*matchedCell, directoryName string) error {
	sort.Slice(validCells, func(i, j int) bool {
		// regex is used before adding, therefore error check can be avoided
		a, _ := strconv.Atoi(validCells[i].CellValue)
		b, _ := strconv.Atoi(validCells[j].CellValue)

		return a < b
	})

	if err := creditOutputCells(sheet, validCells, len(validCells), directoryName); err != nil {
		return err
	}
	return nil
}

func formatSheet(sheet *xlsx.Sheet, out string) error {
	switch out {
	case annual:
		return sheet.SetColWidth(0, 0, 0.67)
	case month:
		return nil
	case shift:
		return nil
	case statement:
		if err := sheet.SetColWidth(0, 0, 11); err != nil {
			return err
		}
		if err := sheet.SetColWidth(2, 2, 15); err != nil {
			return err
		}
		sheet.Rows[5].SetHeight(100)
		return nil
	default:
		return fmt.Errorf("%s is not a valid document type", out)
	}
}

func saveData(sheet *xlsx.Sheet, validCells []*matchedCell, directoryName, out string) error {
	if out == creditCards {
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
			if err.Error() == statementError {
				continue
			}
			return err
		}

		if err := formatSheet(output, out); err != nil {
			return err
		}
		file.Save(getOutputName(directoryName, cell.CellValue, ext))
	}
	return nil
}

func parse(fName, basename, sheetName, regex, out string, colName rune) error {
	if fName == basename {
		fName = fmt.Sprintf("%s.%s", fName, ext)
	}

	xlFile, err := xlsx.OpenFile(fName)
	if err != nil {
		fmt.Printf("unable to open file : %s\nperhaps it does not exist\n", fName)
		return err
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}

	sheet, err := getSheet(xlFile, sheetName)
	if err != nil {
		fmt.Printf("unable to open sheet : %s\nperhaps it does not exist\n", sheetName)
		return err
	}

	index, err := convertRuneToIndex(colName)
	if err != nil {
		return err
	}

	validCells := getValidCells(sheet, r, index)
	directoryName := fmt.Sprintf("%s_output", basename)
	if len(validCells) > 0 {
		if err := createDirIfNotExists(directoryName); err != nil {
			return err
		}
	}

	if err = saveData(sheet, validCells, directoryName, out); err != nil {
		return err
	}
	fmt.Printf("successfully parsed %s\ncheck out %s for files\n", fName, directoryName)
	return nil
}
