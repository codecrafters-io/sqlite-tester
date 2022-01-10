package internal

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"

	_ "embed"

	_ "modernc.org/sqlite"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

//go:embed test_databases/companies.db
var companiesDbContent []byte

var testQueriesForCompanies = []string{
	"SELECT id, name FROM companies WHERE country = 'micronesia'",
	"SELECT id, name FROM companies WHERE country = 'north korea'",
	"SELECT id, name FROM companies WHERE country = 'tonga'",
	"SELECT id, name FROM companies WHERE country = 'eritrea'",
	"SELECT id, name FROM companies WHERE country = 'republic of the congo'",
	"SELECT id, name FROM companies WHERE country = 'montserrat'",
	"SELECT id, name FROM companies WHERE country = 'chad'",
}

func testIndexScan(stageHarness tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	if err := os.WriteFile("./test.db", companiesDbContent, 0666); err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	randomTestQueries := shuffle(testQueriesForCompanies)[0:2]

	for _, testQuery := range randomTestQueries {
		logger.Infof("$ ./your_sqlite3.sh test.db \"%v\"", testQuery)
		result, err := executable.Run("test.db", testQuery)
		if err != nil {
			return err
		}

		actualValues := strings.Split(strings.TrimSpace(string(result.Stdout)), "\n")

		expectedValues, err := getExpectedValuesForQuery(db, testQuery)
		if err != nil {
			logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
			return err
		}

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
	}

	return nil
}
