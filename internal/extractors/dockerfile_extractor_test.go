package extractors

import (
	"testing"

	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
)

func TestExtractImagesFromDockerfiles(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{
		{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile", RelativePath: "Dockerfile"},
		{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile-2", RelativePath: "Dockerfile-2"},
		{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile-3", RelativePath: "Dockerfile-3"},
	}

	envVars := map[string]map[string]string{
		"../../test_files/imageExtraction/dockerfiles": {
			"MY_ARG":     "6.0",
			"MY_ASPNET":  "aspnet",
			"MY_TAG":     "latest",
			"GO_VERSION": "1.20.8",
			"ALPINE_VER": "3.18",
		},
	}

	images, err := ExtractImagesFromDockerfiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	expectedImages := map[string]types.ImageLocation{
		"mcr.microsoft.com/dotnet/sdk:6.0":    {Origin: types.DockerFileOrigin, Path: "Dockerfile"},
		"mcr.microsoft.com/dotnet/aspnet:6.0": {Origin: types.DockerFileOrigin, Path: "Dockerfile"},
		"nginx:latest":                        {Origin: types.DockerFileOrigin, Path: "Dockerfile-2"},
		"mcr.microsoft.com/dotnet/aspnet:4.0": {Origin: types.DockerFileOrigin, Path: "Dockerfile-2"},
		"tonistiigi/xx:1.2.1":                 {Origin: types.DockerFileOrigin, Path: "Dockerfile-3"},
		"golang:1.20.8-alpine3.18":            {Origin: types.DockerFileOrigin, Path: "Dockerfile-3"},
		"alpine:3.18":                         {Origin: types.DockerFileOrigin, Path: "Dockerfile-3"},
	}

	checkResult(t, images, expectedImages)
}

func TestExtractImagesFromDockerfiles_NoFilesFound(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{} // No files provided
	envVars := map[string]map[string]string{}

	images, err := ExtractImagesFromDockerfiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	if len(images) != 0 {
		t.Errorf("Expected 0 images, but got %d", len(images))
	}
}

func TestExtractImagesFromDockerfiles_NoImagesFound(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{
		{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile-4", RelativePath: "Dockerfile-4"}, // Empty Dockerfile
	}

	envVars := map[string]map[string]string{}

	images, err := ExtractImagesFromDockerfiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	if len(images) != 0 {
		t.Errorf("Expected 0 images, but got %d", len(images))
	}
}

func TestExtractImagesFromDockerfiles_OneValidOneInvalid(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{
		{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile", RelativePath: "Dockerfile"},
		{FullPath: "../../test_files/imageExtraction/dockerfiles/InvalidDockerfile", RelativePath: "InvalidDockerfile"},
	}

	envVars := map[string]map[string]string{}

	images, err := ExtractImagesFromDockerfiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	expectedImages := map[string]types.ImageLocation{
		"mcr.microsoft.com/dotnet/sdk:6.0":    {Origin: types.DockerFileOrigin, Path: "Dockerfile"},
		"mcr.microsoft.com/dotnet/aspnet:6.0": {Origin: types.DockerFileOrigin, Path: "Dockerfile"},
	}

	checkResult(t, images, expectedImages)
}

func TestExtractImagesFromDockerfiles_WithEnvFiles(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{
		{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile-5", RelativePath: "Dockerfile-5"},
	}

	envVars := map[string]map[string]string{
		"../../test_files/imageExtraction/dockerfiles": {
			"MY_IMAGE": "golang",
			"MY_TAG":   "1.20.8",
		},
	}

	images, err := ExtractImagesFromDockerfiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	expectedImages := map[string]types.ImageLocation{
		"golang:1.20.8": {Origin: types.DockerFileOrigin, Path: "Dockerfile-5"},
	}

	checkResult(t, images, expectedImages)
}

func TestExtractImagesFromDockerfiles_WithMultipleEnvFiles(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{
		{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile-5", RelativePath: "Dockerfile-5"},
	}

	envVars := map[string]map[string]string{
		"../../test_files/imageExtraction/dockerfiles": {
			"MY_TAG": "3.18",
		},
		"../../test_files/imageExtraction": {
			"MY_IMAGE": "alpine",
		},
	}

	images, err := ExtractImagesFromDockerfiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	expectedImages := map[string]types.ImageLocation{
		"alpine:3.18": {Origin: types.DockerFileOrigin, Path: "Dockerfile-5"},
	}

	checkResult(t, images, expectedImages)
}

func checkResult(t *testing.T, images []types.ImageModel, expectedImages map[string]types.ImageLocation) {
	if len(images) != len(expectedImages) {
		t.Errorf("Expected %d images, but got %d", len(expectedImages), len(images))
	}

	for _, image := range images {
		expectedLocation, ok := expectedImages[image.Name]
		if !ok {
			t.Errorf("Unexpected image found: %s", image.Name)
			continue
		}

		if len(image.ImageLocations) != 1 {
			t.Errorf("Expected image %s to have exactly one location, but got %d", image.Name, len(image.ImageLocations))
			continue
		}

		if image.ImageLocations[0].Path != expectedLocation.Path {
			t.Errorf("Expected image %s to have path %s, but got %s", image.Name, expectedLocation.Path, image.ImageLocations[0].Path)
		}

		if image.ImageLocations[0].Origin != expectedLocation.Origin {
			t.Errorf("Expected image %s to have origin %s, but got %s", image.Name, expectedLocation.Origin, image.ImageLocations[0].Origin)
		}
	}
}
