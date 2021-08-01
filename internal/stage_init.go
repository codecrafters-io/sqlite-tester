package internal

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	_ "modernc.org/sqlite"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testInit(stageHarness tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	tableNames := randomStringsShort(2 + randomInt(7))

	logger.Infof("Creating test database with %v tables: test.db", len(tableNames))

	for _, tableName := range tableNames {
		sqlStmt := fmt.Sprintf(`
			create table %v (id integer not null primary key, name text);
		`, tableName)

		_, err = db.Exec(sqlStmt)
		if err != nil {
			logger.Errorf("Failed to create test table, this is a CodeCrafters error.")
			return err
		}
	}

	logger.Infof("Executing \"./your_sqlite3.sh test.db .dbinfo\"")
	result, err := executable.Run("test.db", ".dbinfo")
	if err != nil {
		return err
	}

	numberOfTablesRegex := regexp.MustCompile(fmt.Sprintf("number of tables:\\s+%v", len(tableNames)))
	numberOfTablesFriendlyPattern := fmt.Sprintf("number of tables: %v", len(tableNames))

	if err = assertStdoutMatchesRegex(result, *numberOfTablesRegex, numberOfTablesFriendlyPattern); err != nil {
		return err
	}

	return nil
}

func assertFileContents(friendlyName string, path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	actualContents := string(bytes)
	expectedContents := "ref: refs/heads/master\n"
	if actualContents != expectedContents {
		return fmt.Errorf("Expected %s to contain '%s', got '%s'", friendlyName, expectedContents, actualContents)
	}

	return nil
}

func assertDirExistsInDir(parent string, child string) error {
	info, err := os.Stat(path.Join(parent, child))
	if os.IsNotExist(err) {
		return fmt.Errorf("Expected the '%s' directory to be created", child)
	}

	if !info.IsDir() {
		return fmt.Errorf("Expected '%s' to be a directory", child)
	}

	return nil
}

func assertFileExistsInDir(parent string, child string) error {
	info, err := os.Stat(path.Join(parent, child))
	if os.IsNotExist(err) {
		return fmt.Errorf("Expected the '%s' file to be created", child)
	}

	if info.IsDir() {
		return fmt.Errorf("Expected '%s' to be a file", child)
	}

	return nil
}

func logDebugTree(logger *tester_utils.Logger, dir string) {
	logger.Debugf("Files found in directory: ")
	doLogDebugTree(logger, dir, " ")
}

func doLogDebugTree(logger *tester_utils.Logger, dir string, prefix string) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	if len(entries) == 0 {
		logger.Debugf(prefix + "  (directory is empty)")
	}

	for _, info := range entries {
		if info.IsDir() {
			logger.Debugf(prefix + "- " + info.Name() + "/")
			doLogDebugTree(logger, path.Join(dir, info.Name()), prefix+" ")
		} else {
			logger.Debugf(prefix + "- " + info.Name())
		}
	}
}
