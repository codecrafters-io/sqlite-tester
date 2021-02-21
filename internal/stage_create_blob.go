package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func testCreateBlob(stageHarness tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	tempDir, err := ioutil.TempDir("", "worktree")
	if err != nil {
		return err
	}

	executable.WorkingDir = tempDir

	logger.Debugf("Running ./your_git.sh init")
	_, err = executable.Run("init")
	if err != nil {
		return err
	}

	logger.Debugf("Writing sample file")
	sampleFile := path.Join(tempDir, fmt.Sprintf("%s.txt", randomStringShort()))
	sampleFileContents := randomString()
	err = ioutil.WriteFile(
		sampleFile,
		[]byte(sampleFileContents),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	expectedSha := plumbing.ComputeHash(plumbing.BlobObject, []byte(sampleFileContents))

	logger.Debugf("Running ./your_git.sh hash-object -w <file>")
	result, err := executable.Run("hash-object", "-w", sampleFile)
	if err != nil {
		return err
	}

	if err = assertExitCode(result, 0); err != nil {
		return err
	}

	if err = assertStdoutContains(result, expectedSha.String()); err != nil {
		return err
	}

	logger.Debugf("Running git cat-file -p %s", expectedSha.String())
	r, err := git.PlainOpen(tempDir)
	if err != nil {
		return err
	}

	blob, err := r.BlobObject(expectedSha)
	if err != nil {
		return err
	}

	blobReader, err := blob.Reader()
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(blobReader)
	if err != nil {
		return err
	}

	expected, actual := sampleFileContents, string(bytes)

	if expected != actual {
		return fmt.Errorf("Expected %q as file contents, got: %q", expected, actual)
	}

	return nil
}
