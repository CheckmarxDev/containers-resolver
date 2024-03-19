package files

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"os"
	"testing"
)

func TestExtractTar(t *testing.T) {
	l := logger.NewLogger(true)

	t.Run("ValidTar", func(t *testing.T) {
		// Provide the path to a valid tar.gz file for testing
		validTarPath := "../../test_files/withDockerInTar.tar.gz"

		extractDir, err := extractTar(l, validTarPath)
		if err != nil {
			t.Fatalf("Error extracting valid tar file: %v", err)
		}
		defer os.RemoveAll(extractDir)

		// Check if the extraction directory exists
		if _, err := os.Stat(extractDir); os.IsNotExist(err) {
			t.Errorf("Extraction directory does not exist: %s", extractDir)
		}
	})

	t.Run("InvalidTar", func(t *testing.T) {
		// Provide the path to an invalid tar.gz file for testing
		invalidTarPath := "../../test_files/invalidWithDockerInTar.tar.gz"

		_, err := extractTar(l, invalidTarPath)
		if err == nil {
			t.Error("Expected error extracting invalid tar file, but got nil")
		}
	})
}
