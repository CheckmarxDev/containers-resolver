package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLogger_Debug(t *testing.T) {
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

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	l := NewLogger(false) // false indicates InfoLevel or higher
	l.Info("Info message")

	expectedOutput := "[INFO] Info message\n"
	actual := buf.String()
	if !strings.Contains(actual, expectedOutput) {
		t.Errorf("Expected info message '%s' to be printed, but it was not. Got: %s", expectedOutput, actual)
	}
}

func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	l := NewLogger(false) // false indicates WarnLevel or higher
	l.Warn("Warn message")

	expectedOutput := "[WARN] Warn message\n"
	actual := buf.String()
	if !strings.Contains(actual, expectedOutput) {
		t.Errorf("Expected warn message '%s' to be printed, but it was not. Got: %s", expectedOutput, actual)
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	l := NewLogger(false) // false indicates ErrorLevel or higher
	l.Error("Error message")

	expectedOutput := "[ERROR] Error message\n"
	actual := buf.String()
	if !strings.Contains(actual, expectedOutput) {
		t.Errorf("Expected error message '%s' to be printed, but it was not. Got: %s", expectedOutput, actual)
	}
}

// Ensure that lower level logs are not printed when a higher level is set
func TestLogger_Level(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	l := NewLogger(false) // false indicates InfoLevel or higher
	l.Debug("This should not print")
	l.Info("This should print")

	expectedOutput := "[INFO] This should print\n"
	actual := buf.String()
	if strings.Contains(actual, "This should not print") {
		t.Errorf("Expected debug message to not be printed, but it was.")
	}
	if !strings.Contains(actual, expectedOutput) {
		t.Errorf("Expected info message '%s' to be printed, but it was not. Got: %s", expectedOutput, actual)
	}
}
