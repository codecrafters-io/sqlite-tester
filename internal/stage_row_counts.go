package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRowCounts(stageHarness *test_case_harness.TestCaseHarness) error {
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
	rowCount := random.RandomInt(100, 200)

	logger.Debugf("Creating table '%v' with %v rows", tableName, rowCount)

	createTableSql := fmt.Sprintf(`create table %v (id integer primary key, name text);`, tableName)

	_, err = db.Exec(createTableSql)
	if err != nil {
		logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
		return err
	}

	insertRowsSql := fmt.Sprintf(
		`insert into %v (name) VALUES %v`,
		tableName,
		strings.Repeat("('dummy_value'), ", rowCount-1)+"('dummy_value')",
	)

	_, err = db.Exec(insertRowsSql)
	if err != nil {
		logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
		return err
	}

	sqlQuery := fmt.Sprintf("select count(*) from %v", tableName)

	logger.Infof("$ ./%v test.db \"%s\"", path.Base(executable.Path), sqlQuery)
	result, err := executable.Run("test.db", sqlQuery)
	if err != nil {
		return err
	}

	if err := assertExitCode(result, 0); err != nil {
		return err
	}

	if err = assertStdout(result, fmt.Sprintf("%v\n", rowCount)); err != nil {
		return err
	}

	return nil
}
