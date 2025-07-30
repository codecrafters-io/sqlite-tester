package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

var expectedQueryResultMap = map[string][]string{
	"SELECT id, name FROM companies WHERE country = 'micronesia'": {
		"1307865|college of micronesia",
		"3696903|nanofabrica",
		"4023193|fsm statistics",
		"6132291|vital energy micronesia",
		"6387751|fsm development bank",
	},
	"SELECT id, name FROM companies WHERE country = 'north korea'": {
		"986681|isn network company limited",
		"1573653|initial innovation limited",
		"2828420|beacon point ltd",
		"3485462|pyongyang university of science & technology (pust)",
		"3969653|plastoform industries ltd",
		"4271599|korea national insurance corporation",
	},
	"SELECT id, name FROM companies WHERE country = 'tonga'": {
		"361142|tonga communications corporation",
		"3186430|tonga development bank",
		"3583436|leiola group limited",
		"4796634|royco amalgamated company limited",
		"7084593|tonga business enterprise centre",
	},
	"SELECT id, name FROM companies WHERE country = 'eritrea'": {
		"121311|unilink s.c.",
		"2102438|orange asmara it solutions",
		"5729848|zara mining share company",
		"6634629|asmara rental",
	},
	"SELECT id, name FROM companies WHERE country = 'republic of the congo'": {
		"509721|skytic telecom",
		"517263|somedia",
		"2543747|its congo",
		"2995059|petroleum trading congo e&p sa",
	},
	"SELECT id, name FROM companies WHERE country = 'montserrat'": {
		"288999|government of montserrat",
		"4472846|university of science, arts & technology",
		"5316703|the abella group llc",
	},
	"SELECT id, name FROM companies WHERE country = 'chad'": {
		"25661|ziyara",
		"987266|hotel la mirande tchad",
		"1313534|kreich avocats",
		"2203192|societe des telecommunications du tchad",
		"2435360|global logistics services limited (gls)",
		"2466228|web tchad",
		"2676248|hanana group",
		"3775391|compagnie sucrière du tchad (cst)",
		"4693857|wenaklabs",
		"5021724|mariam high tech",
		"5255614|bureau d'appui santé et environnement",
		"6828605|tigo tchad",
	},
}

func testIndexScan(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	_ = os.Remove("./test.db")

	if err := os.Symlink(path.Join(os.Getenv("TESTER_DIR"), "companies.db"), "./test.db"); err != nil {
		logger.Errorf("Failed to create symlink for test database, this is a CodeCrafters error.")
		return err
	}

	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		logger.Errorf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	var testQueriesForCompanies []string
	for k := range expectedQueryResultMap {
		testQueriesForCompanies = append(testQueriesForCompanies, k)
	}
	randomTestQueries := random.ShuffleArray(testQueriesForCompanies)[0:2]

	for _, testQuery := range randomTestQueries {
		logger.Infof("$ ./%v test.db \"%v\"", path.Base(executable.Path), testQuery)
		result, err := executable.Run("test.db", testQuery)
		if err != nil {
			return err
		}

		if err := assertExitCode(result, 0); err != nil {
			return err
		}

		actualValues := splitBytesToLines(result.Stdout)

		expectedValues := expectedQueryResultMap[testQuery]
		if len(actualValues) != len(expectedValues) {
			return fmt.Errorf("Expected exactly %v lines of output, got: %v", len(expectedValues), len(actualValues))
		}

		sort.Strings(expectedValues)
		sort.Strings(actualValues)

		expectedValuesStr := strings.Join(expectedValues, "\n")
		actualValuesStr := strings.Join(actualValues, "\n")

		if expectedValuesStr != actualValuesStr {
			return fmt.Errorf("Expected %v to be returned as values, got: %v", expectedValues, actualValues)
		}
	}

	return nil
}
