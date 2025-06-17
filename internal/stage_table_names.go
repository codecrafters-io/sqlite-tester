package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testTableNames(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	tableNames := random.RandomStrings(5)
	sort.Strings(tableNames)

	logger.Debugf("Creating test.db with tables: %v", tableNames)

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

	logger.Infof("$ ./%v test.db .tables", path.Base(executable.Path))
	result, err := executable.Run("test.db", ".tables")
	if err != nil {
		return err
	}

	tableNamesRegex := regexp.MustCompile(fmt.Sprint(strings.Join(tableNames, "\\s+")))
	tableNamesFriendlyPattern := fmt.Sprint(strings.Join(tableNames, " "))

	if err = assertStdoutMatchesRegex(result, *tableNamesRegex, tableNamesFriendlyPattern); err != nil {
		return err
	}

	return nil
}
