.PHONY: release build test test_with_git copy_course_file

current_version_number := $(shell git tag --list "v*" | sort -V | tail -n 1 | cut -c 2-)
next_version_number := $(shell echo $$(($(current_version_number)+1)))

release:
	git tag v$(next_version_number)
	git push origin master v$(next_version_number)

build:
	go build -o dist/main.out ./cmd/tester

test:
	TESTER_DIR=$(shell pwd) go test -v ./internal/

test_and_watch:
	onchange '**/*' -- go test -v ./internal/

test_with_sqlite: build
	CODECRAFTERS_SUBMISSION_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{"slug":"init","tester_log_prefix":"stage-1","title":"Stage #1: Print page size"},{"slug":"table_count","tester_log_prefix":"stage-2","title":"Stage #2: Print number of tables"},{"slug":"table_names","tester_log_prefix":"stage-3","title":"Stage #3: Print table names"},{"slug":"row_counts","tester_log_prefix":"stage-4","title":"Stage #4: Count rows in a table"},{"slug":"read_single_column","tester_log_prefix":"stage-5","title":"Stage #5: Read data from a single column"},{"slug":"read_multiple_columns","tester_log_prefix":"stage-6","title":"Stage #6: Read data from multiple columns"},{"slug":"where","tester_log_prefix":"stage-7","title":"Stage #7: Filter data with a WHERE clause"},{"slug":"table_scan","tester_log_prefix":"stage-8","title":"Stage #8: Retrieve data using a full-table scan"},{"slug":"index_scan","tester_log_prefix":"stage-9","title":"Stage #9: Retrieve data using an index"}]" \
	dist/main.out

copy_course_file:
	hub api \
		repos/rohitpaulk/codecrafters-server/contents/codecrafters/store/data/sqlite.yml \
		| jq -r .content \
		| base64 -d \
		> internal/test_helpers/course_definition.yml

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils
