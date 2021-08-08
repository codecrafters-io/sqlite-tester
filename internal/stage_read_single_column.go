package internal

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bxcodec/faker/v3"

	_ "modernc.org/sqlite"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

const NUMBER_OF_COLUMNS = 5

func testReadSingleColumn(stageHarness tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	tableName := randomStringShort()
	allColumnNames := randomStringsShort(NUMBER_OF_COLUMNS)
	testColumnIndex := randomInt(NUMBER_OF_COLUMNS)
	testColumnName := allColumnNames[testColumnIndex]
	numberOfRecords := 4 + randomInt(4)

	logger.Debugf("Creating test.db with table: %v", tableName)

	columnWithTypeList := []string{}

	for _, columnName := range allColumnNames {
		columnWithTypeList = append(columnWithTypeList, fmt.Sprintf("%v text", columnName))
	}

	createTableSql := fmt.Sprintf(`
      create table %v (id primary key, %v);
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

	logger.Infof("$ ./your_sqlite3.sh test.db \"select %v from %v\"", testColumnName, tableName)
	result, err := executable.Run("test.db", fmt.Sprintf("select %v from %v", testColumnName, tableName))
	if err != nil {
		return err
	}

	actualValues := strings.Split(strings.TrimSpace(string(result.Stdout)), "\n")

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
