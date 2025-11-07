package testutils

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

// AssertJSONContains checks if JSON output contains expected key-value pairs
func AssertJSONContains(t *testing.T, output string, expected map[string]interface{}) {
	t.Helper()
	
	var data interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
		return
	}

	for key, expectedValue := range expected {
		if !containsKeyValue(data, key, expectedValue) {
			t.Errorf("Expected JSON to contain '%s': '%v', but it doesn't", key, expectedValue)
		}
	}
}

// containsKeyValue recursively searches for a key-value pair in JSON data
func containsKeyValue(data interface{}, key string, value interface{}) bool {
	switch v := data.(type) {
	case map[string]interface{}:
		if val, ok := v[key]; ok {
			return val == value
		}
		// Recursively search nested objects
		for _, nestedVal := range v {
			if containsKeyValue(nestedVal, key, value) {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if containsKeyValue(item, key, value) {
				return true
			}
		}
	}
	return false
}

// AssertOutputMatches checks if output matches a regex pattern
func AssertOutputMatches(t *testing.T, output, pattern string) {
	t.Helper()
	matched, err := regexp.MatchString(pattern, output)
	if err != nil {
		t.Errorf("Invalid regex pattern: %v", err)
		return
	}
	if !matched {
		t.Errorf("Expected output to match pattern '%s', but got: %s", pattern, output)
	}
}

// AssertOutputLineCount checks if output has expected number of lines
func AssertOutputLineCount(t *testing.T, output string, expectedLines int) {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	actualLines := len(lines)
	if actualLines != expectedLines {
		t.Errorf("Expected %d lines, but got %d", expectedLines, actualLines)
	}
}

// AssertOutputStartsWith checks if output starts with expected text
func AssertOutputStartsWith(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.HasPrefix(strings.TrimSpace(output), expected) {
		t.Errorf("Expected output to start with '%s', but got: %s", expected, output)
	}
}

// AssertOutputEndsWith checks if output ends with expected text
func AssertOutputEndsWith(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.HasSuffix(strings.TrimSpace(output), expected) {
		t.Errorf("Expected output to end with '%s', but got: %s", expected, output)
	}
}

