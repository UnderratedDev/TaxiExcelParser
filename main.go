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
	return str == Annual || str == Month || str == Shift || str == Statement || str == CreditCards
}

const (
	Annual      string = "annual"
	Month       string = "month"
	Shift       string = "shift"
	Statement   string = "statement"
	CreditCards string = "credit"
	// Extra April
	AnnualReportColName rune   = 'E'
	AnnualReportRegex   string = "^\\d+[a-zA-Z]+ Total$"

	// 4 April
	MonthReportColName  rune   = 'B'
	MonthReportRegex    string = "^\\d+[a-zA-Z]+$"
	MonthReportFGColour string = "FFFF00"

	// Shift Report
	ShiftReportColName rune   = 'D'
	ShiftReportRegex   string = "^\\d+(AM|PM) Total$"

	// Statement
	// Note: Perhaps check 1A numbers instead, however the id is required
	// StatementChunk int = 52
	// StatementErrorSize int = 44
	// StatementColName   rune = 'A'
	// StatementRegex     string = "^\\d+$"
	StatementColName rune   = 'B'
	StatementRegex   string = "^\\d+[a-zA-Z]$"

	CreditColName rune   = 'O'
	CreditRegex   string = "^\\d+$"

	Exit string = "quit"
	Ext  string = "xlsx"
)

type MatchedCell struct {
	RowIndex  int
	ColIndex  int
	CellValue string
}

func main() {
	for {
		out, err := readInputWithPrompt(fmt.Sprintf("To parse a report, enter %s, %s, %s, %s, or %s (type %s to exit) :", Annual, Month, Shift, Statement, CreditCards, Exit))
		if err != nil {
			panic(err)
		}
		if out == Exit {
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
				case Annual:
					parse(fName, basename, sheetName, AnnualReportRegex, out, AnnualReportColName)
					// if err := parse(fName, sheetName, AnnualReportRegex, out, AnnualReportColName); err != nil {
					// }
				case Month:
					parse(fName, basename, sheetName, MonthReportRegex, out, MonthReportColName)
					// if err := parse(fName, sheetName, MonthReportRegex, out, MonthReportColName); err != nil {
					// }
				case Shift:
					parse(fName, basename, sheetName, ShiftReportRegex, out, ShiftReportColName)
					// if err := parse(fName, sheetName, ShiftReportRegex, out, ShiftReportColName); err != nil {
					// }
				case Statement:
					parse(fName, basename, sheetName, StatementRegex, out, StatementColName)
					// if err := parse(fName, sheetName, StatementRegex, out, StatementColName); err != nil {
					// }
				case CreditCards:
					parse(fName, basename, sheetName, CreditRegex, out, CreditColName)
					// if err := parse(fName, sheetName, CreditRegex, out, CreditColName); err != nil {
					// }
				}
			}()
		} else {
			// panic(fmt.Errorf("%s is not a valid type", out))
			// fmt.Printf("%s is not a vaild command, use %s, %s, %s, %s, or %s (type %s to exit) :\n", out, Annual, Month, Shift, Statement, CreditCards, Exit)
			fmt.Printf("%s is not a vaild command\n", out)
		}
	}
}
