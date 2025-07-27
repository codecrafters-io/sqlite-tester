package internal

import (
	"os"
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	os.Setenv("CODECRAFTERS_RANDOM_SEED", "1234567890")

	testCases := map[string]tester_utils_testing.TesterOutputTestCase{
		"init_failure": {
			UntilStageSlug:      "dr6",
			CodePath:            "./test_helpers/stages/init_failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_success": {
			UntilStageSlug:      "dr6",
			CodePath:            "./test_helpers/stages/init",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"table_count_failure": {
			UntilStageSlug:      "ce0",
			CodePath:            "./test_helpers/stages/init",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/table_count/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"table_count_success": {
			UntilStageSlug:      "ce0",
			CodePath:            "./test_helpers/stages/table_count",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/table_count/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"table_names_failure": {
			UntilStageSlug:      "sz4",
			CodePath:            "./test_helpers/stages/table_count",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/table_names/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"table_names_success": {
			UntilStageSlug:      "sz4",
			CodePath:            "./test_helpers/stages/table_names",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/table_names/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"pass_all": {
			// omitted "dr6", "ce0" because sqlite3 on linux is not compiled to support .dbinfo
			// omitted "az9", "vc9", "rf3" because of randomness issues
			StageSlugs:          []string{"sz4", "nd9", "ws9", "nz8"},
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/base/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"segfault": {
			UntilStageSlug:      "ws9",
			CodePath:            "./test_helpers/scenarios/segfault",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/scenarios/segfault",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	return testerOutput
}

// TODO: TestReadSingleColumnFailure
