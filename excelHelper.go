package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

func getSheet(xlFile *xlsx.File, name string) (*xlsx.Sheet, error) {
	for i := range xlFile.Sheets {
		if xlFile.Sheets[i].Name == name {
			return xlFile.Sheets[i], nil
		}
	}
	return nil, fmt.Errorf("sheet %s does not exist in excel file", name)
}

func copyCell(a *xlsx.Cell, b *xlsx.Cell) error {
	value, err := a.FormattedValue()
	if err != nil {
		return err
	}

	b.SetString(value)
	b.SetStyle(a.GetStyle())

	return nil
}

func copyCellWithFill(a *xlsx.Cell, b *xlsx.Cell, fgColour string) error {
	str, err := a.FormattedValue()
	if err != nil {
		return err
	}

	b.SetString(str)
	style := a.GetStyle()
	if len(style.Fill.FgColor) == 0 {
		style.Fill.FgColor = fgColour
		style.ApplyFill = true
	}

	b.SetStyle(style)
	return nil
}

func copyCells(sheet, output *xlsx.Sheet, validCells []*matchedCell, index int) error {
	outputRow := output.AddRow()
	for _, cell := range sheet.Rows[index].Cells {
		outputCell := outputRow.AddCell()
		if err := copyCell(cell, outputCell); err != nil {
			return err
		}
	}

	return nil
}

func xlsxToCsv(name, sheetName string) error {
	xlFile, err := xlsx.OpenFile(fmt.Sprintf("%s.%s", name, ext))
	if err != nil {
		return err
	}

	sheet, err := getSheet(xlFile, sheetName)
	if err != nil {
		return err
	}

	csv, err := os.Create(fmt.Sprintf("%s.csv", name))
	if err != nil {
		return err
	}

	for _, row := range sheet.Rows {
		var line []string
		if row != nil {
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil || len(str) == 0 {
					continue
				}
				line = append(line, strings.TrimSpace(str))
			}
			csv.WriteString(strings.Join(line, ",") + "\n")
		}
	}

	return nil
}
