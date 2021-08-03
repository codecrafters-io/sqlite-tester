package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatStages:    []testerutils.Stage{},
	ExecutableFileName: "your_sqlite3.sh",
	Stages: []testerutils.Stage{
		{
			Slug:                    "init",
			Title:                   "Print number of tables",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Slug:                    "table_names",
			Title:                   "Print table names",
			TestFunc:                testTableNames,
			ShouldRunPreviousStages: true,
		},
		{
			Slug:                    "row_counts",
			Title:                   "Count rows in a table",
			TestFunc:                testRowCounts,
			ShouldRunPreviousStages: true,
		},
		{
			Slug:                    "read_single_column",
			Title:                   "Read data from a single column",
			TestFunc:                testReadSingleColumn,
			ShouldRunPreviousStages: true,
		},
	},
}
