package internal

import (
	"time"

	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatStages:    []testerutils.Stage{},
	ExecutableFileName: "your_sqlite3.sh",
	Stages: []testerutils.Stage{
		{
			Number:                  1,
			Slug:                    "init",
			Title:                   "Print page size",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  2,
			Slug:                    "table_count",
			Title:                   "Print number of tables",
			TestFunc:                testTableCount,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  3,
			Slug:                    "table_names",
			Title:                   "Print table names",
			TestFunc:                testTableNames,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  4,
			Slug:                    "row_counts",
			Title:                   "Count rows in a table",
			TestFunc:                testRowCounts,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  5,
			Slug:                    "read_single_column",
			Title:                   "Read data from a single column",
			TestFunc:                testReadSingleColumn,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  6,
			Slug:                    "read_multiple_columns",
			Title:                   "Read data from multiple columns",
			TestFunc:                testReadMultipleColumns,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  7,
			Slug:                    "where",
			Title:                   "Filter data with a WHERE clause",
			TestFunc:                testWhere,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  8,
			Slug:                    "table_scan",
			Title:                   "Retrieve data using a full-table scan",
			TestFunc:                testTableScan,
			ShouldRunPreviousStages: true,
			Timeout:                 60 * time.Second, // TODO: Turn this back down once we're able to figure out why running inside firecracker takes so long
		},
		{
			Number:                  9,
			Slug:                    "index_scan",
			Title:                   "Retrieve data using an index",
			TestFunc:                testIndexScan,
			ShouldRunPreviousStages: true,
			Timeout:                 20 * time.Second,
		},
	},
}
