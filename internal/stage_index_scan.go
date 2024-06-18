package internal

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

var testQueriesForCompanies = []string{
	"SELECT id, name FROM companies WHERE country = 'micronesia'",
	"SELECT id, name FROM companies WHERE country = 'north korea'",
	"SELECT id, name FROM companies WHERE country = 'tonga'",
	"SELECT id, name FROM companies WHERE country = 'eritrea'",
	"SELECT id, name FROM companies WHERE country = 'republic of the congo'",
	"SELECT id, name FROM companies WHERE country = 'montserrat'",
	"SELECT id, name FROM companies WHERE country = 'chad'",
}

func testIndexScan(stageHarness *test_case_harness.TestCaseHarness) error {
	initRandom()

	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	if err := exec.Command("cp", path.Join(os.Getenv("TESTER_DIR"), "companies.db"), "./test.db").Run(); err != nil {
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
		logger.Infof("$ ./%v test.db \"%v\"", path.Base(executable.Path), testQuery)
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
