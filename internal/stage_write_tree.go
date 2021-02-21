package internal

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	tester_utils "github.com/codecrafters-io/tester-utils"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/filesystem/dotgit"
)

func testWriteTree(stageHarness tester_utils.StageHarness) error {
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

	logger.Debugf("Creating some files & directories")

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

	r, err := git.PlainOpen(tempDir)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil
	}

	// If we're running against git, update the index
	if err = w.AddGlob("."); err != nil {
		return err
	}

	logger.Debugf("Running ./your_git.sh write-tree")
	result, err := executable.Run("write-tree")
	if err != nil {
		return err
	}

	if err = assertExitCode(result, 0); err != nil {
		return err
	}

	sha := strings.TrimSpace(string(result.Stdout))
	if len(sha) != 40 {
		return fmt.Errorf("Expected a 40-char SHA as output. Got: %v", sha)
	}

	logger.Debugf("Running git ls-tree --name-only <sha>")

	storage := filesystem.NewObjectStorage(
		dotgit.New(osfs.New(path.Join(tempDir, ".git"))),
		cache.NewObjectLRU(0),
	)
	tree, err := object.GetTree(storage, plumbing.NewHash(sha))
	if err != nil {
		return err
	}

	actual := ""
	for _, entry := range tree.Entries {
		actual += entry.Name
		actual += "\n"
	}

	expectedValues := []string{rootFile, rootDir1, rootDir2}
	sort.Strings(expectedValues)
	expected := strings.Join(
		expectedValues,
		"\n",
	) + "\n"

	if expected != actual {
		return fmt.Errorf("Expected %q as stdout, got: %q", expected, actual)
	}

	return nil
}
