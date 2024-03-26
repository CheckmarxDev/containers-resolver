package files

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFindHelmCharts(t *testing.T) {
	baseDir := "../../test_files/imageExtraction"

	// Test the findHelmCharts function
	helmCharts, err := findHelmCharts(baseDir)
	if err != nil {
		t.Fatalf("Error finding Helm charts: %v", err)
	}

	// Expected Helm chart info
	expectedChart := types.HelmChartInfo{
		Directory:  filepath.Join(baseDir, "helm"),
		ValuesFile: "helm/values.yaml",
		TemplateFiles: []types.FilePath{
			{FullPath: "../../test_files/imageExtraction/helm/templates/containers-worker.yaml", RelativePath: "helm/templates/containers-worker.yaml"},
			{FullPath: "../../test_files/imageExtraction/helm/templates/image-insights.yaml", RelativePath: "helm/templates/image-insights.yaml"},
		},
	}

	// Verify if the retrieved Helm chart info matches the expected
	if len(helmCharts) != 1 {
		t.Fatalf("Expected 1 Helm chart, got %d", len(helmCharts))
	}
	if !reflect.DeepEqual(helmCharts[0], expectedChart) {
		t.Errorf("Retrieved Helm chart info does not match expected:\nGot: %+v\nExpected: %+v", helmCharts[0], expectedChart)
	}
}

func TestIsValidFolderPath(t *testing.T) {
	baseDir := "../../test_files/imageExtraction"

	// Test with a valid folder path
	isValid, err := IsValidFolderPath(baseDir)
	if err != nil {
		t.Errorf("Error while checking valid folder path: %v", err)
	}
	if !isValid {
		t.Errorf("Expected valid folder path, got invalid")
	}

	// Test with an invalid folder path
	isValid, err = IsValidFolderPath(baseDir + "/nonexistent")
	if err == nil || isValid {
		t.Errorf("Expected invalid folder path, got valid")
	}
}

func TestDeleteDirectory(t *testing.T) {
	tempDir := setupTempDir(t)
	defer cleanupTempDir(t, tempDir)

	filePath := tempDir + "/testfile.txt"
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	file.Close()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Test file does not exist before deletion")
	}

	err = DeleteDirectory(tempDir)
	if err != nil {
		t.Errorf("Error deleting directory: %v", err)
	}

	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Errorf("Directory still exists after deletion")
	}
}

func TestExtractCompressedPath_Zip(t *testing.T) {
	baseDir := "../../test_files/withDockerInZip.zip"

	tempDir := setupTempDir(t)
	defer cleanupTempDir(t, tempDir)

	extractedPath, err := extractCompressedPath(logger.NewLogger(false), baseDir)
	if err != nil {
		t.Fatalf("Error extracting zip file: %v", err)
	}

	_, err = os.Stat(filepath.Join(extractedPath, "Dockerfile"))
	if err != nil {
		t.Errorf("Expected Dockerfile in extracted directory, got error: %v", err)
	}
}

func TestExtractCompressedPath_Tar(t *testing.T) {
	baseDir := "../../test_files/withDockerInTar.tar.gz"

	tempDir := setupTempDir(t)
	defer cleanupTempDir(t, tempDir)

	extractedPath, err := extractCompressedPath(logger.NewLogger(false), baseDir)
	if err != nil {
		t.Fatalf("Error extracting tar.gz file: %v", err)
	}

	_, err = os.Stat(filepath.Join(extractedPath, "withDockerInTar/Dockerfile"))
	if err != nil {
		t.Errorf("Expected Dockerfile in extracted directory, got error: %v", err)
	}
}

func TestExtractCompressedPath_NonExistingFile(t *testing.T) {
	nonExistingFile := "nonexistingfile.zip"

	_, err := extractCompressedPath(logger.NewLogger(false), nonExistingFile)
	if err == nil {
		t.Errorf("Expected error for non-existing file, got nil")
	}
}
func TestGetRelativePath(t *testing.T) {
	// Test case for relative path within the same directory
	baseDir := "/path/to/base/dir"
	filePath := "/path/to/base/dir/file.txt"
	expectedRelativePath := "file.txt"
	if relativePath := getRelativePath(baseDir, filePath); relativePath != expectedRelativePath {
		t.Errorf("Expected relative path %s, got %s", expectedRelativePath, relativePath)
	}

	// Test case for relative path in a subdirectory
	baseDir = "/path/to/base"
	filePath = "/path/to/base/dir/file.txt"
	expectedRelativePath = "dir/file.txt"
	if relativePath := getRelativePath(baseDir, filePath); relativePath != expectedRelativePath {
		t.Errorf("Expected relative path %s, got %s", expectedRelativePath, relativePath)
	}
}

func TestIsYAMLFile(t *testing.T) {
	// Test case for YAML file
	yamlFilePath := "/path/to/file.yaml"
	if !isYAMLFile(yamlFilePath) {
		t.Errorf("Expected %s to be recognized as a YAML file", yamlFilePath)
	}

	// Test case for non-YAML file
	nonYAMLFilePath := "/path/to/file.txt"
	if isYAMLFile(nonYAMLFilePath) {
		t.Errorf("Expected %s not to be recognized as a YAML file", nonYAMLFilePath)
	}
}

func TestIsHelmChart(t *testing.T) {
	// Test case for Helm chart directory
	helmChartDir := "../../test_files/imageExtraction/helm"
	if !isHelmChart(helmChartDir) {
		t.Errorf("Expected %s to be recognized as a Helm chart directory", helmChartDir)
	}

	// Test case for non-Helm chart directory
	nonHelmChartDir := "/path/to/non/helm/chart"
	if isHelmChart(nonHelmChartDir) {
		t.Errorf("Expected %s not to be recognized as a Helm chart directory", nonHelmChartDir)
	}
}

func TestGetContainerResolutionFullPath(t *testing.T) {
	// Test case for getting container resolution full path
	folderPath := "/path/to/folder"
	expectedFullPath := "/path/to/folder/containers-resolution.json"
	if fullPath, _ := getContainerResolutionFullPath(folderPath); fullPath != expectedFullPath {
		t.Errorf("Expected full path %s, got %s", expectedFullPath, fullPath)
	}
}

func setupTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Error creating temp directory: %v", err)
	}
	return dir
}

func cleanupTempDir(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Error cleaning up temp directory: %v", err)
	}
}
