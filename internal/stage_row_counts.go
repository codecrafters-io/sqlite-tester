package internal

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "modernc.org/sqlite"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testRowCounts(stageHarness tester_utils.StageHarness) error {
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
	rowCount := randomInt(200)

	logger.Debugf("Creating table '%v' with %v rows", tableName, rowCount)

	createTableSql := fmt.Sprintf(`create table %v (id integer not null primary key, name text);`, tableName)

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

	logger.Infof("$ ./your_sqlite3.sh test.db \"select count(*) from test1\"")
	result, err := executable.Run("test.db", fmt.Sprintf("select count(*) from %v", tableName))
	if err != nil {
		return err
	}

	if err = assertStdout(result, fmt.Sprintf("%v\n", rowCount)); err != nil {
		return err
	}

	return nil
}
