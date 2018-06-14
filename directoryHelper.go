package main

import (
	"fmt"
	"os"
)

func createDirIfNotExists(directoryName string) error {
	if _, err := os.Stat(directoryName); os.IsNotExist(err) {
		if err = os.Mkdir(directoryName, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func checkIfExists(path string) error {
	if _, err := os.Stat(path); os.IsExist(err) {
		return err
	}
	return nil
}

func getOutputName(directoryName, cellValue, ext string) string {
	outputName := fmt.Sprintf("%s/%s.%s", directoryName, cellValue, ext)
	for fileIndex, err := 0, checkIfExists(outputName); err != nil; fileIndex++ {
		outputName = fmt.Sprintf("%s/%s_%d.%s", directoryName, cellValue, fileIndex, Ext)
	}

	return outputName
}
