package internal

import (
	"time"

	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatTestCases:    []testerutils.TestCase{},
	ExecutableFileName: "your_sqlite3.sh",
	TestCases: []testerutils.TestCase{
		{
			Slug:                    "init",
			TestFunc:                testInit,
		},
		{
			Slug:                    "table_count",
			TestFunc:                testTableCount,
		},
		{
			Slug:                    "table_names",
			TestFunc:                testTableNames,
		},
		{
			Slug:                    "row_counts",
			TestFunc:                testRowCounts,
		},
		{
			Slug:                    "read_single_column",
			TestFunc:                testReadSingleColumn,
		},
		{
			Slug:                    "read_multiple_columns",
			TestFunc:                testReadMultipleColumns,
		},
		{
			Slug:                    "where",
			TestFunc:                testWhere,
		},
		{
			Slug:                    "table_scan",
			TestFunc:                testTableScan,
			Timeout:                 60 * time.Second, // TODO: Turn this back down once we're able to figure out why running inside firecracker takes so long
		},
		{
			Slug:                    "index_scan",
			TestFunc:                testIndexScan,
			Timeout:                 20 * time.Second,
		},
	},
}
