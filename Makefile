SHELL := /bin/bash
.PHONY: default build clean test lint dist release do_dist set_version_from_file set_version_from_git dev

export GOPATH := $(GOPATH):$(shell pwd):$(shell pwd)/vendor
default: build

build:
	go get github.com/constabulary/gb/...
	go version
	gb build

clean:
	@rm -Rf bin pkg .build build dist test

test_pkg: clean build
	@RES=0; \
	COUNT=$$((COUNT+1)); \
	go test -tags test -race $$pkg; \
	if [ $$? -ne 0 ]; then \
		RES=1; \
	fi; \
	exit $$RES

test: clean build lint
	@RES=0; \
	TEST_OUTPUT_DIR=test; \
	COVERAGE_FILE=$$TEST_OUTPUT_DIR"/coverage.out"; \
	EXCLUDED_PKG_LIST=(\
		"dolo-tracking" \
		"dolo-tracking/context" \
		"dolo-tracking/logger" \
		); \
	mkdir -p $$TEST_OUTPUT_DIR; \
	echo "mode: count" > $$COVERAGE_FILE; \
	FAILURE_MESSAGES="[Error]: Some of the tests failed. Run these commands to pinpoint the cause:\n"; \
	COUNT=0; \
	for folder in $(shell find src -type d); do \
		COUNT=$$((COUNT+1)); \
		if [ "$$folder" != "src" ]; then \
			file_count=`find $$folder/ -maxdepth 1 -name *.go | wc -l`; \
			if [ $$file_count -gt 0 ]; then \
				pkg=`echo $$folder | sed 's,^[^/]*/,,'`; \
				if [[ ! " $${EXCLUDED_PKG_LIST[@]} " =~ " $${pkg} " ]]; then \
					go test -tags test -coverprofile=$$TEST_OUTPUT_DIR/$$COUNT-coverage.out -race $$pkg; \
					if [ $$? -ne 0 ]; then \
						FAILURE_MESSAGES=$$FAILURE_MESSAGES"\n\t make test_pkg pkg=$$pkg"; \
						RES=1; \
					elif [ ! -f $$TEST_OUTPUT_DIR/$$COUNT-coverage.out ]; then \
						FAILURE_MESSAGES=$$FAILURE_MESSAGES"\n\tNo tests for package: $$pkg"; \
						RES=1; \
					else \
						tail -n +2 $$TEST_OUTPUT_DIR/$$COUNT-coverage.out >> $$COVERAGE_FILE; \
					fi; \
					rm -f $$TEST_OUTPUT_DIR/$$COUNT-coverage.out; \
				fi \
			fi \
		fi \
	done; \
	if [ $$RES -ne 0 ]; then \
		echo ""; \
		echo -e $$FAILURE_MESSAGES; \
		echo ""; \
		exit $$RES; \
	fi; \
	go tool cover -html=$$COVERAGE_FILE -o $$TEST_OUTPUT_DIR/coverage.html; \
	go tool cover -func=$$COVERAGE_FILE -o $$TEST_OUTPUT_DIR/coverage.txt; \
	echo "("`tail -1 test/coverage.txt  | tr -d '[:space:]' | cut -d")" -f2`") covered"; \
	rm -f $$COVERAGE_FILE; \
	exit 0

lint:
	@RES=0; \
	FAILURE_MESSAGES="[Error]: Some of the checks failed. Run these commands to pinpoint the cause:\n"; \
	go get -u github.com/golang/lint/golint; \
	for folder in $(shell find src -type d); do \
		if [ "$$folder" != "src" ]; then \
			file_count=`find $$folder/ -maxdepth 1 -name *.go | wc -l`; \
			if [ $$file_count -gt 0 ]; then \
				pkg=`echo $$folder | sed 's,^[^/]*/,,'`; \
				echo "Linting and vetting go package: $$pkg"; \
				golint -set_exit_status $$pkg; \
				if [ $$? -ne 0 ]; then \
					FAILURE_MESSAGES=$$FAILURE_MESSAGES"\n\tgolint $$pkg"; \
					RES=1; \
				fi; \
				go vet $$pkg; \
				if [ $$? -ne 0 ]; then \
					FAILURE_MESSAGES=$$FAILURE_MESSAGES"\n\tgo vet $$pkg"; \
					RES=1; \
				fi; \
			fi \
		fi \
	done; \
	if [ $$RES -ne 0 ]; then \
		echo ""; \
		echo -e $$FAILURE_MESSAGES; \
		echo ""; \
		exit $$RES; \
	fi; \
	exit 0
