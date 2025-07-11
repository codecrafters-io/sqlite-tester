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

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

var testQueriesForSuperheroes = []string{
	"SELECT id, name FROM superheroes WHERE eye_color = 'Pink Eyes'",
	"SELECT id, name FROM superheroes WHERE eye_color = 'Amber Eyes'",
	"SELECT id, name FROM superheroes WHERE eye_color = 'Gold Eyes'",
	"SELECT id, name FROM superheroes WHERE hair_color = 'Reddish Brown Hair'",
	"SELECT id, name FROM superheroes WHERE hair_color = 'Silver Hair'",
	"SELECT id, name FROM superheroes WHERE hair_color = 'Violet Hair'",
	"SELECT id, name FROM superheroes WHERE hair_color = 'Gold Hair'",
}

func testTableScan(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	if err := exec.Command("cp", path.Join(os.Getenv("TESTER_DIR"), "superheroes.db"), "./test.db").Run(); err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	randomTestQueries := random.ShuffleArray(testQueriesForSuperheroes)[0:3]

	for _, testQuery := range randomTestQueries {
		logger.Infof("$ ./%v test.db \"%v\"", path.Base(executable.Path), testQuery)
		result, err := executable.Run("test.db", testQuery)
		if err != nil {
			return err
		}

		if err := assertExitCode(result, 0); err != nil {
			return err
		}

		actualValues := splitBytesToLines(result.Stdout)

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

func getExpectedValuesForQuery(db *sql.DB, query string) ([]string, error) {
	expectedValues := []string{}

	rows, err := db.Query(query)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var value1 string
		var value2 string

		if err := rows.Scan(&value1, &value2); err != nil {
			return []string{}, err
		}

		expectedValues = append(expectedValues, strings.Join([]string{value1, value2}, "|"))
	}

	return expectedValues, nil
}
