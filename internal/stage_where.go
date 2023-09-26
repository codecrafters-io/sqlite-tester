package internal

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"

	_ "modernc.org/sqlite"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testWhere(stageHarness *tester_utils.StageHarness) error {
	initRandom()

	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	table := generateRandomTable()
	numberOfRecords := 4 + randomInt(4)

	logger.Debugf("Creating test.db with table: %v", table.Name)

	_, err = db.Exec(table.CreateTableSQL())
	if err != nil {
		logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
		return err
	}

	records := []Record{}

	for i := 1; i <= numberOfRecords; i++ {
		records = append(records, generateRandomRecord(table))
	}

	_, err = db.Exec(table.InsertRecordsSQL(records))
	if err != nil {
		logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
		return err
	}

	testColumnNames := shuffle(table.ColumnNames)[0:3]
	expectedValues := []string{}

	filterColumnName := shuffle(table.ColumnNames)[0]
	filterColumnValue := records[0].ValueFor(filterColumnName)

	for _, record := range records {
		if record.ValueFor(filterColumnName) == filterColumnValue {
			expectedValues = append(expectedValues, strings.Join(record.ValuesFor(testColumnNames), "|"))
		}
	}

	selectColumnsSql := fmt.Sprintf(
		"select %v from %v where %v = '%v'",
		strings.Join(testColumnNames, ", "),
		table.Name,
		filterColumnName,
		filterColumnValue,
	)

	logger.Infof("$ ./your_sqlite3.sh test.db \"%v\"", selectColumnsSql)
	result, err := executable.Run("test.db", selectColumnsSql)
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
