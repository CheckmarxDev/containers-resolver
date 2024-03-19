package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLogger_DebugMode(t *testing.T) {
	testCases := []struct {
		debugModeEnabled bool
		message          string
		expectedOutput   string
	}{
		{debugModeEnabled: true, message: "Debug message", expectedOutput: "[DEBUG] Debug message\n"},
		{debugModeEnabled: false, message: "Debug message", expectedOutput: ""},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		log.SetOutput(&buf)

		l := NewLogger(tc.debugModeEnabled)
		l.Debug(tc.message)

		actual := buf.String()
		if tc.debugModeEnabled && !strings.Contains(actual, tc.expectedOutput) {
			t.Errorf("Expected debug message '%s' to be printed in debug mode, but it was not. Got: %s", tc.expectedOutput, actual)
		} else if !tc.debugModeEnabled && actual != "" {
			t.Errorf("Expected no debug message to be printed when debug mode is off, but got: %s", actual)
		}
		log.SetOutput(nil)
	}
}
