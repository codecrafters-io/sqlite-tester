.PHONY: release build test test_with_git copy_course_file

current_version_number := $(shell git tag --list "v*" | sort -V | tail -n 1 | cut -c 2-)
next_version_number := $(shell echo $$(($(current_version_number)+1)))

release:
	git tag v$(next_version_number)
	git push origin master v$(next_version_number)

build:
	go build -o dist/main.out ./cmd/tester
	go build -o dist/starter.out ./cmd/starter_tester

test:
	go test -v ./internal/

test_and_watch:
	onchange '**/*' -- go test -v ./internal/

test_with_sqlite: build
	CODECRAFTERS_SUBMISSION_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_CURRENT_STAGE_SLUG="row_counts" \
	CODECRAFTERS_COURSE_PAGE_URL="test" \
	dist/main.out

copy_course_file:
	hub api \
		repos/rohitpaulk/codecrafters-server/contents/codecrafters/store/data/sqlite.yml \
		| jq -r .content \
		| base64 -d \
		> internal/test_helpers/course_definition.yml
