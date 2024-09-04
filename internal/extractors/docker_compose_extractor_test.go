package extractors

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"testing"
)

func TestExtractImagesFromDockerComposeFiles(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{
		{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose.yaml", RelativePath: "docker-compose.yaml"},
		{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose-2.yaml", RelativePath: "docker-compose-2.yaml"},
	}

	envVars := map[string]map[string]string{}

	images, err := ExtractImagesFromDockerComposeFiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	expectedImages := map[string]types.ImageLocation{
		"postgres:12.0": {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-2.yaml"},
		"minio/minio:RELEASE.2020-06-22T03-12-50Z": {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-2.yaml"},
		"redis:6.0.10-alpine":                      {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-2.yaml"},
		"source.azure.io/api:latest":               {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-2.yaml"},
		"buildimage:latest":                        {Origin: types.DockerComposeFileOrigin, Path: "docker-compose.yaml"},
	}

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

func TestExtractImagesFromDockerComposeFile(t *testing.T) {
	l := logger.NewLogger(false)
	filePath := types.FilePath{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose-2.yaml", RelativePath: "docker-compose-2.yaml"}

	envVars := map[string]map[string]string{
		"../../test_files/imageExtraction/dockerCompose": {
			"MARKETER_IMAGE": "source.azure.io/api:3.18",
		},
	}

	images, err := extractImagesFromDockerComposeFile(l, filePath, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	expectedImages := map[string]types.ImageLocation{
		"postgres:12.0": {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-2.yaml"},
		"minio/minio:RELEASE.2020-06-22T03-12-50Z": {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-2.yaml"},
		"redis:6.0.10-alpine":                      {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-2.yaml"},
		"source.azure.io/api:3.18":                 {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-2.yaml"},
	}

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

func TestExtractImagesFromDockerComposeFiles_NoFilesFound(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{} // No files provided

	envVars := map[string]map[string]string{}

	images, err := ExtractImagesFromDockerComposeFiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	if len(images) != 0 {
		t.Errorf("Expected 0 images, but got %d", len(images))
	}
}

func TestExtractImagesFromDockerComposeFiles_NoImagesFound(t *testing.T) {
	l := logger.NewLogger(false)

	filePaths := []types.FilePath{
		{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose-5.yaml", RelativePath: "docker-compose-5.yaml"}, // Empty Docker Compose file
	}

	envVars := map[string]map[string]string{}

	images, err := ExtractImagesFromDockerComposeFiles(l, filePaths, envVars)
	if err != nil {
		t.Errorf("Error extracting images: %v", err)
	}

	if len(images) != 0 {
		t.Errorf("Expected 0 images, but got %d", len(images))
	}
}
