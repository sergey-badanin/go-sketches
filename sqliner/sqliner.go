package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
)

/*
 * Small util performs transformation of multi-line SQL query to one line
 * Useful when query is to be applied to postgresql \copy command
 *
 * Program accepts 2 command line arguments:
 * 1st: file name for source SQL
 * 2nd: file name for resulting sql
 */
func main() {
	sourceFileName, resultFileName := getFileNames()
	translateToOneLine(sourceFileName, resultFileName)
}

func getFileNames() (sourceFileName string, resultFileName string) {
	if len(os.Args) < 3 {
		log.Fatal("It seems some parameters are not provided")
	}
	sourceFileName = os.Args[1]
	resultFileName = os.Args[2]
	return
}

func translateToOneLine(sourceFileName string, resultFileName string) {
	sourceFile, err := os.Open(sourceFileName)
	check(err)
	defer sourceFile.Close()

	resultFile, err := os.Create(resultFileName)
	check(err)
	defer resultFile.Close()

	resultString := loadFileToString(sourceFile)
	replacer := strings.NewReplacer("\n", "", "\r", "")
	resultString = replacer.Replace(resultString)

	replacer = strings.NewReplacer("  ", " ")
	for i := strings.Index(resultString, "  "); i != -1; i = strings.Index(resultString, "  ") {
		resultString = replacer.Replace(resultString)
	}

	replacer = strings.NewReplacer("( ", "(", " )", ")")
	resultString = replacer.Replace(resultString)

	resultString = strings.TrimSpace(resultString)

	fileWriter := bufio.NewWriter(resultFile)
	fileWriter.WriteString(resultString)
	fileWriter.Flush()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func loadFileToString(file *os.File) string {
	sqlScanner := bufio.NewScanner(file)

	var resultBuffer bytes.Buffer
	for sqlScanner.Scan() {
		resultBuffer.WriteString(" " + sqlScanner.Text())
	}
	check(sqlScanner.Err())
	return resultBuffer.String()
}
