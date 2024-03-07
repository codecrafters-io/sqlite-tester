package internal

import (
	"time"

	"github.com/codecrafters-io/tester-utils/tester_definition"
)

var testerDefinition = tester_definition.TesterDefinition{
	AntiCheatTestCases: []tester_definition.TestCase{},
	ExecutableFileName: "your_sqlite3.sh",
	TestCases: []tester_definition.TestCase{
		{
			Slug:     "init",
			TestFunc: testInit,
		},
		{
			Slug:     "table_count",
			TestFunc: testTableCount,
		},
		{
			Slug:     "table_names",
			TestFunc: testTableNames,
		},
		{
			Slug:     "row_counts",
			TestFunc: testRowCounts,
		},
		{
			Slug:     "read_single_column",
			TestFunc: testReadSingleColumn,
		},
		{
			Slug:     "read_multiple_columns",
			TestFunc: testReadMultipleColumns,
		},
		{
			Slug:     "where",
			TestFunc: testWhere,
		},
		{
			Slug:     "table_scan",
			TestFunc: testTableScan,
			Timeout:  60 * time.Second, // TODO: Turn this back down once we're able to figure out why running inside firecracker takes so long
		},
		{
			Slug:     "index_scan",
			TestFunc: testIndexScan,
			Timeout:  20 * time.Second,
		},
	},
}
