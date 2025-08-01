package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"regexp"

	_ "modernc.org/sqlite"

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testTableCount(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	tableNames := random.RandomWords(2 + random.RandomInt(0, 7))

	logger.Debugf("Creating test database with %v tables: test.db", len(tableNames))

	for _, tableName := range tableNames {
		sqlStmt := fmt.Sprintf(`
			create table %v (id integer primary key, name text);
		`, tableName)

		_, err = db.Exec(sqlStmt)
		if err != nil {
			logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
			return err
		}
	}

	logger.Infof("$ ./%v test.db .dbinfo", path.Base(executable.Path))
	result, err := executable.Run("test.db", ".dbinfo")
	if err != nil {
		return err
	}

	if err := assertExitCode(result, 0); err != nil {
		return err
	}

	numberOfTablesRegex := regexp.MustCompile(fmt.Sprintf("number of tables:\\s+%v", len(tableNames)))
	numberOfTablesFriendlyPattern := fmt.Sprintf("number of tables: %v", len(tableNames))

	if err = assertStdoutMatchesRegex(result, *numberOfTablesRegex, numberOfTablesFriendlyPattern); err != nil {
		return err
	}

	return nil
}
