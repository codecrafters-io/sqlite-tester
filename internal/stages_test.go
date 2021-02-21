package internal

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("init", "test_helpers/stages/init_failure")
	if !assert.Equal(t, 1, exitCode) {
		failWithMockerOutput(t, m)
	}
	assert.Contains(t, m.ReadStdout(), "nothing")
	assert.Contains(t, m.ReadStdout(), "Test failed")

	m.Reset()

	exitCode = runCLIStage("init", "test_helpers/stages/init")
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}
}

func TestReadBlob(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("read_blob", "test_helpers/stages/init")
	if !assert.Equal(t, 1, exitCode) {
		failWithMockerOutput(t, m)
	}
	if !assert.Contains(t, m.ReadStdout(), "Expected") {
		failWithMockerOutput(t, m)
	}
	if !assert.Contains(t, m.ReadStdout(), "as stdout") {
		failWithMockerOutput(t, m)
	}
	if !assert.Contains(t, m.ReadStdout(), "Test failed") {
		failWithMockerOutput(t, m)
	}

	m.Reset()

	exitCode = runCLIStage("read_blob", "test_helpers/stages/read_blob")
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}
}

func TestCreateBlob(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("create_blob", "test_helpers/stages/read_blob")
	if !assert.Equal(t, 1, exitCode) {
		failWithMockerOutput(t, m)
	}
	assert.Contains(t, m.ReadStdout(), "Expected")
	assert.Contains(t, m.ReadStdout(), "Test failed")

	m.Reset()

	exitCode = runCLIStage("create_blob", "test_helpers/stages/create_blob")
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}
}

func TestReadTree(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("read_tree", "test_helpers/stages/read_blob")
	if !assert.Equal(t, 1, exitCode) {
		failWithMockerOutput(t, m)
	}
	assert.Contains(t, m.ReadStdout(), "Expected")
	assert.Contains(t, m.ReadStdout(), "Test failed")

	m.Reset()

	exitCode = runCLIStage("read_tree", "test_helpers/stages/read_tree")
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}
}

func TestWriteTree(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("write_tree", "test_helpers/stages/read_tree")
	if !assert.Equal(t, 1, exitCode) {
		failWithMockerOutput(t, m)
	}
	assert.Contains(t, m.ReadStdout(), "Expected")
	assert.Contains(t, m.ReadStdout(), "Test failed")

	m.Reset()

	exitCode = runCLIStage("write_tree", "test_helpers/stages/write_tree")
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}
}

func TestCreateCommit(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("create_commit", "test_helpers/stages/write_tree")
	if !assert.Equal(t, 1, exitCode) {
		failWithMockerOutput(t, m)
	}
	assert.Contains(t, m.ReadStdout(), "Expected")
	assert.Contains(t, m.ReadStdout(), "Test failed")

	m.Reset()

	t.Skip("Not done yet")

	exitCode = runCLIStage("create_commit", "test_helpers/stages/create_commit")
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}
}

func runCLIStage(slug string, dir string) (exitCode int) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return RunCLI(map[string]string{
		"CODECRAFTERS_CURRENT_STAGE_SLUG": slug,
		"CODECRAFTERS_SUBMISSION_DIR":     path.Join(cwd, dir),
	})
}

func failWithMockerOutput(t *testing.T, m *IOMocker) {
	m.End()
	t.Error(fmt.Sprintf("stdout: \n%s\n\nstderr: \n%s", m.ReadStdout(), m.ReadStderr()))
	t.FailNow()
}
