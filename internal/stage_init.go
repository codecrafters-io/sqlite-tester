package internal

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"

	_ "modernc.org/sqlite"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testInit(stageHarness *tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	pageSizes := []int{512, 1024, 2048, 4096, 8192, 16384, 32768}
	pageSize := pageSizes[randomInt(len(pageSizes))]

	logger.Debugf("Creating test database with page size %d: test.db", pageSize)
	db, err := sql.Open("sqlite", fmt.Sprintf("./test.db?_pragma=page_size(%d)", pageSize))
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE test (id integer primary key, name text);")
	if err != nil {
		logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
		return err
	}

	logger.Infof("$ ./your_sqlite3.sh test.db .dbinfo")
	result, err := executable.Run("test.db", ".dbinfo")
	if err != nil {
		return err
	}

	databasePageSizeRegex := regexp.MustCompile(fmt.Sprintf("database page size:\\s+%v", pageSize))
	databasePageSizeFriendlyPattern := fmt.Sprintf("database page size: %v", pageSize)

	if err = assertStdoutMatchesRegex(result, *databasePageSizeRegex, databasePageSizeFriendlyPattern); err != nil {
		return err
	}

	return nil
}
