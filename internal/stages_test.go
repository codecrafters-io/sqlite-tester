package internal

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	testCases := map[string]tester_utils_testing.TesterOutputTestCase{
		"init_failure": {
			UntilStageSlug: "init",
			CodePath: "./test_helpers/stages/init_failure",
			ExpectedExitCode: 1,
			StdoutFixturePath: "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_success": {
			UntilStageSlug: "init",
			CodePath: "./test_helpers/stages/init",
			ExpectedExitCode: 0,
			StdoutFixturePath: "./test_helpers/fixtures/init/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"table_count_failure": {
			UntilStageSlug: "table_count",
			CodePath: "./test_helpers/stages/init",
			ExpectedExitCode: 1,
			StdoutFixturePath: "./test_helpers/fixtures/table_count/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"table_count_success": {
			UntilStageSlug: "table_count",
			CodePath: "./test_helpers/stages/table_count",
			ExpectedExitCode: 0,
			StdoutFixturePath: "./test_helpers/fixtures/table_count/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"table_names_failure": {
			UntilStageSlug: "table_names",
			CodePath: "./test_helpers/stages/table_count",
			ExpectedExitCode: 1,
			StdoutFixturePath: "./test_helpers/fixtures/table_names/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"table_names_success": {
			UntilStageSlug: "table_names",
			CodePath: "./test_helpers/stages/table_names",
			ExpectedExitCode: 0,
			StdoutFixturePath: "./test_helpers/fixtures/table_names/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	return testerOutput
}

// TODO: TestReadSingleColumnFailure

