package internal

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	tester_utils "github.com/codecrafters-io/tester-utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func testCreateCommit(stageHarness tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	tempDir, err := ioutil.TempDir("", "worktree")
	if err != nil {
		return err
	}

	executable.WorkingDir = tempDir

	logger.Debugf("Running git init")
	if _, err := git.PlainInit(tempDir, false); err != nil {
		return err
	}

	logger.Debugf("Creating some files & directories")

	rootFile := "root.txt"
	firstLevel := randomStringsShort(3)
	rootFile, rootDir1, rootDir2 := firstLevel[0], firstLevel[1], firstLevel[2]
	secondLevel := randomStringsShort(2)
	rootDir1File1, rootDir1File2 := secondLevel[0], secondLevel[1]
	thirdLevel := randomStringsShort(2)
	rootDir2File1, rootDir2File2 := thirdLevel[0], thirdLevel[1]

	writeFile(tempDir, rootFile)
	writeFile(tempDir, path.Join(rootDir1, rootDir1File1))
	writeFile(tempDir, path.Join(rootDir1, rootDir1File2))
	writeFile(tempDir, path.Join(rootDir2, rootDir2File1))

	logger.Debugf("Running git commit --all")
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

	commitHash, err := w.Commit("test2", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	parentCommitSha := commitHash.String()

	logger.Debugf("Creating another file")
	writeFile(tempDir, path.Join(rootDir2, rootDir2File2))

	if err = w.AddGlob("."); err != nil {
		return err
	}

	nextCommitHash, err := w.Commit("test2", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	nextCommit, err := repository.CommitObject(nextCommitHash)
	if err != nil {
		return err
	}

	treeSha := nextCommit.TreeHash.String()

	commitMessage := randomString()
	logger.Debugf("Running ./your_git.sh commit-tree <tree_sha> -p <commit_sha> -m <message>")
	result, err := executable.Run("commit-tree", treeSha, "-p", parentCommitSha, "-m", commitMessage)
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

	logger.Debugf("Running git cat-file commit <sha>")
	c, err := repository.CommitObject(plumbing.NewHash(sha))
	if err != nil {
		return err
	}

	expected, actual := treeSha, c.TreeHash.String()
	if expected != actual {
		return fmt.Errorf("Expected %q as tree, got: %q", expected, actual)
	}

	parent, err := c.Parent(0)
	if err != nil {
		return err
	}
	expected, actual = parentCommitSha, parent.Hash.String()
	if expected != actual {
		return fmt.Errorf("Expected %q as parent commit, got: %q", expected, actual)
	}

	expected, actual = commitMessage+"\n", c.Message
	if expected != actual {
		return fmt.Errorf("Expected %q as commit message, got: %q", expected, actual)
	}

	return nil
}
