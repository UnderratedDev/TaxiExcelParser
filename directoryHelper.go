package main

import (
	"fmt"
	"os"
)

func createDirIfNotExists(directoryName string) error {
	if err := checkIfExists(directoryName); err == nil {
		return os.Mkdir(directoryName, os.ModePerm)
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
		outputName = fmt.Sprintf("%s/%s_%d.%s", directoryName, cellValue, fileIndex, ext)
	}

	return outputName
}
