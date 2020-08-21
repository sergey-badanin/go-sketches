package main

import (
	"os"
	"testing"
)

func TestTranslateToOneLine(t *testing.T) {
	sourceFileName := "test.sql"
	sourceFile, err := os.Open(sourceFileName)
	check(err)
	defer sourceFile.Close()
	sourceSql := loadFileToString(sourceFile)

	expectedFile, err := os.Open("expected.sql")
	check(err)
	defer expectedFile.Close()
	expectedSql := loadFileToString(expectedFile)

	resultingFileName := "resulting.sql"
	translateToOneLine(sourceFileName, resultingFileName)

	resultingFile, err := os.Open(resultingFileName)
	check(err)
	defer resultingFile.Close()
	resultingSql := loadFileToString(resultingFile)

	if resultingSql != expectedSql {
		t.Errorf("For \n%v\n got \n%v\n, expected \n%v", sourceSql, resultingSql, expectedSql)
	}
}
