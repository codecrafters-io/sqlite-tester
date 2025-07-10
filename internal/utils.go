package internal

import "strings"

func splitBytesToLines(input []byte) []string {
	trimmed := strings.TrimSpace(string(input))
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "\n")
}
