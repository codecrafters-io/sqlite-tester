package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatStages:    []testerutils.Stage{},
	ExecutableFileName: "your_git.sh",
	Stages: []testerutils.Stage{
		{
			Slug:     "init",
			Title:    "Initialize the .git directory",
			TestFunc: testInit,
		},
		{
			Slug:     "read_blob",
			Title:    "Read a blob object",
			TestFunc: testReadBlob,
		},
		{
			Slug:     "create_blob",
			Title:    "Create a blob object",
			TestFunc: testCreateBlob,
		},
		{
			Slug:     "read_tree",
			Title:    "Read a tree object",
			TestFunc: testReadTree,
		},
		{
			Slug:     "write_tree",
			Title:    "Write a tree object",
			TestFunc: testWriteTree,
		},
		{
			Slug:     "create_commit",
			Title:    "Create a commit object",
			TestFunc: testCreateCommit,
		},
		{
			Slug:     "clone_repository",
			Title:    "Clone a repository",
			TestFunc: testCloneRepository,
		},
	},
}
