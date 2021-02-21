package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	tester_utils "github.com/codecrafters-io/tester-utils"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/filesystem/dotgit"
)

func testReadTree(stageHarness tester_utils.StageHarness) error {
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

	logger.Debugf("Writing a tree to git storage..")

	rootFile := "root.txt"
	firstLevel := randomStringsShort(3)
	rootFile, rootDir1, rootDir2 := firstLevel[0], firstLevel[1], firstLevel[2]
	secondLevel := randomStringsShort(2)
	rootDir1File1, rootDir1File2 := secondLevel[0], secondLevel[1]
	rootDir2File1 := randomStringShort()

	writeFile(tempDir, rootFile)
	writeFile(tempDir, path.Join(rootDir1, rootDir1File1))
	writeFile(tempDir, path.Join(rootDir1, rootDir1File2))
	writeFile(tempDir, path.Join(rootDir2, rootDir2File1))

	repository, err := git.PlainOpen(tempDir)
	if err != nil {
		return err
	}

	w, err := repository.Worktree()
	if err != nil {
		return err
	}

	if err = w.AddGlob("."); err != nil {
		return err
	}

	commitHash, err := w.Commit("test", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	commit, err := repository.CommitObject(commitHash)
	if err != nil {
		return err
	}

	sha := commit.TreeHash.String()
	logger.Debugf("Running ./your_git.sh ls-tree --name-only %s", sha)
	result, err := executable.Run("ls-tree", "--name-only", sha)
	if err != nil {
		return err
	}

	storage := filesystem.NewObjectStorage(
		dotgit.New(osfs.New(path.Join(tempDir, ".git"))),
		cache.NewObjectLRU(0),
	)
	tree, err := object.GetTree(storage, plumbing.NewHash(sha))
	if err != nil {
		return err
	}

	expected := ""

	for _, entry := range tree.Entries {
		expected += entry.Name
		expected += "\n"
	}

	actual := string(result.Stdout)

	if expected != actual {
		return fmt.Errorf("Expected %q as stdout, got: %q", expected, actual)
	}

	return nil
}

func writeFile(rootDir string, filepath string) {
	filepath = path.Join(rootDir, filepath)
	if err := os.MkdirAll(path.Dir(filepath), 0700); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(filepath, []byte(randomString()), os.ModePerm); err != nil {
		panic(err)
	}
}
