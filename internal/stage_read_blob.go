package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	tester_utils "github.com/codecrafters-io/tester-utils"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/filesystem/dotgit"
)

func testReadBlob(stageHarness tester_utils.StageHarness) error {
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

	storage := filesystem.NewObjectStorage(
		dotgit.New(osfs.New(path.Join(tempDir, ".git"))),
		cache.NewObjectLRU(0),
	)
	obj := storage.NewEncodedObject()
	obj.SetType(plumbing.BlobObject)
	writer, err := obj.Writer()
	if err != nil {
		return err
	}

	if _, err := writer.Write([]byte(sampleFileContents)); err != nil {
		return err
	}

	hash, err := storage.SetEncodedObject(obj)
	if err != nil {
		return err
	}

	if hash != expectedSha {
		panic("Expected sha doesn't match!")
	}

	logger.Debugf("Running ./your_git.sh cat-file -p %s", expectedSha.String())
	result, err := executable.Run("cat-file", "-p", expectedSha.String())
	if err != nil {
		return err
	}

	if err = assertExitCode(result, 0); err != nil {
		return err
	}

	return assertStdout(result, sampleFileContents)
}
