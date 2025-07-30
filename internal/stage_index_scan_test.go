package internal

import (
	"database/sql"
	"fmt"
	"sort"
	"testing"

	_ "modernc.org/sqlite"
)

func TestExpectedQueryResultMap(t *testing.T) {
	db, err := sql.Open("sqlite", "../companies.db")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	for query, expectedResults := range expectedQueryResultMap {
		t.Run(fmt.Sprintf("query_%s", query), func(t *testing.T) {
			// Execute the query
			rows, err := db.Query(query)
			if err != nil {
				t.Fatalf("Failed to execute query '%s': %v", query, err)
			}
			defer rows.Close()

			var actualResults []string
			for rows.Next() {
				var id int
				var name string
				if err := rows.Scan(&id, &name); err != nil {
					t.Fatalf("Failed to scan row: %v", err)
				}
				actualResults = append(actualResults, fmt.Sprintf("%d|%s", id, name))
			}

			if err := rows.Err(); err != nil {
				t.Fatalf("Error iterating rows: %v", err)
			}

			sort.Strings(expectedResults)
			sort.Strings(actualResults)

			if len(actualResults) != len(expectedResults) {
				t.Errorf("Expected %d results, got %d results", len(expectedResults), len(actualResults))
				t.Errorf("Expected: %v", expectedResults)
				t.Errorf("Actual: %v", actualResults)
				return
			}

			for i, expected := range expectedResults {
				if i >= len(actualResults) {
					t.Errorf("Missing result at index %d: expected %s", i, expected)
					continue
				}
				if actualResults[i] != expected {
					t.Errorf("Result mismatch at index %d: expected %s, got %s", i, expected, actualResults[i])
				}
			}

			if t.Failed() {
				t.Logf("Query: %s", query)
				t.Logf("Expected results: %v", expectedResults)
				t.Logf("Actual results: %v", actualResults)
			}
		})
	}
}
