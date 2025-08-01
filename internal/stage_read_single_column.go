package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/bxcodec/faker/v3"

	_ "modernc.org/sqlite"

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

const NUMBER_OF_COLUMNS = 5

func testReadSingleColumn(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	tableName := random.RandomWord()
	allColumnNames := random.RandomWords(NUMBER_OF_COLUMNS)
	testColumnIndex := random.RandomInt(0, NUMBER_OF_COLUMNS)
	testColumnName := allColumnNames[testColumnIndex]
	numberOfRecords := 4 + random.RandomInt(0, 4)

	logger.Debugf("Creating test.db with table: %v", tableName)
	logger.Debugf("Columns in table: %s", strings.Join(allColumnNames, ", "))

	columnWithTypeList := []string{}

	for _, columnName := range allColumnNames {
		columnWithTypeList = append(columnWithTypeList, fmt.Sprintf("%v text", columnName))
	}

	createTableSql := fmt.Sprintf(`
      create table %v (id integer primary key, %v);
    `, tableName, strings.Join(columnWithTypeList, ","))

	_, err = db.Exec(createTableSql)
	if err != nil {
		logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
		return err
	}

	generateValuesForRecord := func() []string {
		values := []string{}

		for i := 1; i <= NUMBER_OF_COLUMNS; i++ {
			values = append(values, faker.FirstNameFemale())
		}

		return values
	}

	recordValuesList := [][]string{}

	for i := 1; i <= numberOfRecords; i++ {
		recordValuesList = append(recordValuesList, generateValuesForRecord())
	}

	expectedValues := []string{}
	valuesSqlList := []string{}

	for _, recordValues := range recordValuesList {
		expectedValues = append(expectedValues, recordValues[testColumnIndex])
		valuesSqlList = append(valuesSqlList, "('"+strings.Join(recordValues, "' , '")+"')")
	}

	insertRowsSql := fmt.Sprintf(
		`insert into %v (%v) VALUES %v`,
		tableName,
		strings.Join(allColumnNames, ", "),
		strings.Join(valuesSqlList, ", "),
	)

	_, err = db.Exec(insertRowsSql)
	if err != nil {
		logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
		return err
	}

	logger.Infof("$ ./%v test.db \"select %v from %v\"", path.Base(executable.Path), testColumnName, tableName)
	result, err := executable.Run("test.db", fmt.Sprintf("select %v from %v", testColumnName, tableName))
	if err != nil {
		return err
	}

	if err := assertExitCode(result, 0); err != nil {
		return err
	}

	actualValues := splitBytesToLines(result.Stdout)

	if len(actualValues) != len(expectedValues) {
		return fmt.Errorf("Expected exactly %v lines of output, got: %v", len(expectedValues), len(actualValues))
	}

	sort.Strings(expectedValues)
	sort.Strings(actualValues)

	expectedValuesStr := strings.Join(expectedValues, "\n")
	actualValuesStr := strings.Join(actualValues, "\n")

	if expectedValuesStr != actualValuesStr {
		return fmt.Errorf("Expected %v to be returned as values, got: %v", expectedValues, actualValues)
	}

	return nil
}
