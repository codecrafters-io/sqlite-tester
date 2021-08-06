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

//go:embed test_databases/superheroes.db
var superheroesDbContent []byte

var testQueries = []string{
	"SELECT id, name FROM superheroes WHERE eye_color = 'Pink Eyes'",
	"SELECT id, name FROM superheroes WHERE eye_color = 'Amber Eyes'",
	"SELECT id, name FROM superheroes WHERE eye_color = 'Gold Eyes'",
	"SELECT id, name FROM superheroes WHERE hair_color = 'Reddish Brown Hair'",
	"SELECT id, name FROM superheroes WHERE hair_color = 'Silver Hair'",
	"SELECT id, name FROM superheroes WHERE hair_color = 'Violet Hair'",
	"SELECT id, name FROM superheroes WHERE hair_color = 'Gold Hair'",
}

func testTableScan(stageHarness tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	if err := os.WriteFile("./test.db", superheroesDbContent, 0666); err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	randomTestQueries := shuffle(testQueries)[0:3]

	for _, testQuery := range randomTestQueries {
		logger.Infof("Executing: ./your_sqlite3.sh test.db \"%v\"", testQuery)
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
