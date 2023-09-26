package internal

import (
	"fmt"
	"regexp"
	"strings"

	executable "github.com/codecrafters-io/tester-utils/executable"
)

func assertEqual(actual string, expected string) error {
	if expected != actual {
		return fmt.Errorf("Expected %q as stdout, got: %q", expected, actual)
	}

	return nil
}

func assertStdout(result executable.ExecutableResult, expected string) error {
	actual := string(result.Stdout)
	if expected != actual {
		return fmt.Errorf("Expected %q as stdout, got: %q", expected, actual)
	}

	return nil
}

func assertStderr(result executable.ExecutableResult, expected string) error {
	actual := string(result.Stderr)
	if expected != actual {
		return fmt.Errorf("Expected %q as stderr, got: %q", expected, actual)
	}

	return nil
}

func assertStdoutContains(result executable.ExecutableResult, expectedSubstring string) error {
	actual := string(result.Stdout)
	if !strings.Contains(actual, expectedSubstring) {
		return fmt.Errorf("Expected stdout to contain %q, got: %q", expectedSubstring, actual)
	}

	return nil
}

func assertStdoutMatchesRegex(result executable.ExecutableResult, pattern regexp.Regexp, friendlyPattern string) error {
	actual := string(result.Stdout)
	if !pattern.MatchString(actual) {
		return fmt.Errorf("Expected stdout to contain %q, got: %q", friendlyPattern, actual)
	}

	return nil
}

func assertStderrContains(result executable.ExecutableResult, expectedSubstring string) error {
	actual := string(result.Stderr)
	if !strings.Contains(actual, expectedSubstring) {
		return fmt.Errorf("Expected stderr to contain %q, got: %q", expectedSubstring, actual)
	}

	return nil
}

func assertExitCode(result executable.ExecutableResult, expected int) error {
	actual := result.ExitCode
	if expected != actual {
		return fmt.Errorf("Expected %d as exit code, got: %d", expected, actual)
	}

	return nil
}
