package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testInit(stageHarness tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	logger.Debugf("Running git init")
	tempDir, err := ioutil.TempDir("", "worktree")
	if err != nil {
		return err
	}

	executable.WorkingDir = tempDir
	_, err = executable.Run("init")
	if err != nil {
		return err
	}

	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err = assertDirExistsInDir(tempDir, dir); err != nil {
			logDebugTree(logger, tempDir)
			return err
		}
	}

	for _, file := range []string{".git/HEAD"} {
		if err = assertFileExistsInDir(tempDir, file); err != nil {
			logDebugTree(logger, tempDir)
			return err
		}
	}

	if err = assertFileContents(".git/HEAD", path.Join(tempDir, ".git/HEAD")); err != nil {
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
