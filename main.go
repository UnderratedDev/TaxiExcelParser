package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readInputWithPrompt(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func isValidType(typ string) bool {
	str := strings.TrimSpace(typ)
	return str == annual || str == month || str == shift || str == statement || str == creditCards
}

const (
	// Types of supported documents
	annual      string = "annual"
	month       string = "month"
	shift       string = "shift"
	statement   string = "statement"
	creditCards string = "credit"
	// Extra April
	annualReportColName rune   = 'E'
	annualReportRegex   string = "^\\d+[a-zA-Z]+ Total$"

	// 4 April
	monthReportColName  rune   = 'B'
	monthReportRegex    string = "^\\d+[a-zA-Z]+$"
	monthReportFGColour string = "FFFF00"

	// Shift Report
	shiftReportColName rune   = 'B'
	shiftReportRegex   string = "^\\d+(Day|Night) Total$"

	// Statement
	statementColName rune   = 'B'
	statementRegex   string = "^\\d+[a-zA-Z]$"

	// Credit
	creditColName rune   = 'K'
	creditRegex   string = "^\\d+$"

	exit string = "quit"
	ext  string = "xlsx"
)

type matchedCell struct {
	RowIndex  int
	ColIndex  int
	CellValue string
}

func main() {
	for {
		out, err := readInputWithPrompt(fmt.Sprintf("To parse a report, enter %s, %s, %s, %s, or %s (type %s to exit) : ", annual, month, shift, statement, creditCards, exit))
		if err != nil {
			panic(err)
		}
		if out == exit {
			os.Exit(1)
		} else if isValidType(out) {
			fName, sheetName, err := getNames()
			if err != nil {
				panic(err)
			}

			basename := getBaseName(fName)
			go func() {
				fmt.Printf("working on %s in the background\n", fName)
				switch out {
				case annual:
					parse(fName, basename, sheetName, annualReportRegex, out, annualReportColName)
				case month:
					parse(fName, basename, sheetName, monthReportRegex, out, monthReportColName)
				case shift:
					parse(fName, basename, sheetName, shiftReportRegex, out, shiftReportColName)
				case statement:
					parse(fName, basename, sheetName, statementRegex, out, statementColName)
				case creditCards:
					parse(fName, basename, sheetName, creditRegex, out, creditColName)
				}
			}()
		} else {
			fmt.Printf("%s is not a vaild command\n", out)
		}
	}
}
