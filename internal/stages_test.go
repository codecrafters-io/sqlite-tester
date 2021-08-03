package internal

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFailure(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("init", "test_helpers/stages/init_failure")
	if !assert.Equal(t, 1, exitCode) {
		failWithMockerOutput(t, m)
	}
	assert.Contains(t, m.ReadStdout(), "nothing")
	assert.Contains(t, m.ReadStdout(), "number of tables")
	assert.Contains(t, m.ReadStdout(), "Test failed")
}

func TestInitSuccess(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("init", "test_helpers/stages/init")
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}
}

func TestTableNamesFailure(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("table_names", "test_helpers/stages/init")
	if !assert.Equal(t, 1, exitCode) {
		failWithMockerOutput(t, m)
	}

	m.End()

	assert.Contains(t, m.ReadStdout(), "Invalid command")
	assert.Contains(t, m.ReadStdout(), "Expected stdout to contain")
	assert.Contains(t, m.ReadStdout(), "Test failed")
}

func TestTableNamesSuccess(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("table_names", "test_helpers/stages/table_names")
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}
}

func TestReadSingleColumnFailure(t *testing.T) {
	m := NewStdIOMocker()
	m.Start()
	defer m.End()

	exitCode := runCLIStage("read_single_column", "test_helpers/pass_all") // TODO: Change to row counts
	if !assert.Equal(t, 0, exitCode) {
		failWithMockerOutput(t, m)
	}

	failWithMockerOutput(t, m)

	////m.End()
	//
	//assert.Contains(t, m.ReadStdout(), "Invalid command")
	//assert.Contains(t, m.ReadStdout(), "Expected stdout to contain")
	//assert.Contains(t, m.ReadStdout(), "Test failed")
}

func runCLIStage(slug string, dir string) (exitCode int) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return RunCLI(map[string]string{
		"CODECRAFTERS_CURRENT_STAGE_SLUG": slug,
		"CODECRAFTERS_SUBMISSION_DIR":     path.Join(cwd, dir),
		"CODECRAFTERS_COURSE_PAGE_URL":    "test",
	})
}

func failWithMockerOutput(t *testing.T, m *IOMocker) {
	m.End()
	t.Error(fmt.Sprintf("stdout: \n%s\n\nstderr: \n%s", m.ReadStdout(), m.ReadStderr()))
	t.FailNow()
}
