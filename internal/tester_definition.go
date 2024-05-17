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
			Slug:     "dr6",
			TestFunc: testInit,
		},
		{
			Slug:     "ce0",
			TestFunc: testTableCount,
		},
		{
			Slug:     "sz4",
			TestFunc: testTableNames,
		},
		{
			Slug:     "nd9",
			TestFunc: testRowCounts,
		},
		{
			Slug:     "az9",
			TestFunc: testReadSingleColumn,
		},
		{
			Slug:     "vc9",
			TestFunc: testReadMultipleColumns,
		},
		{
			Slug:     "rf3",
			TestFunc: testWhere,
		},
		{
			Slug:     "ws9",
			TestFunc: testTableScan,
			Timeout:  60 * time.Second, // TODO: Turn this back down once we're able to figure out why running inside firecracker takes so long
		},
		{
			Slug:     "nz8",
			TestFunc: testIndexScan,
			Timeout:  20 * time.Second,
		},
	},
}
