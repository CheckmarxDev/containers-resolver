package files

import (
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"reflect"
	"testing"
)

func TestMergeImages(t *testing.T) {
	// Define sample image data
	dockerfileImages := []types.ImageModel{
		{Name: "nginx:latest", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
	}
	dockerComposeImages := []types.ImageModel{
		{Name: "mysql:5.7", ImageLocations: []types.ImageLocation{{Origin: types.DockerComposeFileOrigin, Path: "docker-compose.yaml"}}},
	}
	helmImages := []types.ImageModel{
		{Name: "rabbitmq:3", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "values.yaml"}}},
	}

	// Call the mergeImages function
	mergedImages := mergeImages(nil, dockerfileImages, dockerComposeImages, helmImages)

	// Define the expected result
	expectedResult := []types.ImageModel{
		{Name: "nginx:latest", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
		{Name: "mysql:5.7", ImageLocations: []types.ImageLocation{{Origin: types.DockerComposeFileOrigin, Path: "docker-compose.yaml"}}},
		{Name: "rabbitmq:3", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "values.yaml"}}},
	}

	// Check if the merged images match the expected result
	if !reflect.DeepEqual(mergedImages, expectedResult) {
		t.Errorf("MergeImages result does not match the expected result")
	}
}

func TestMergeDuplicates(t *testing.T) {
	// Define sample image data with duplicates
	imageModels := []types.ImageModel{
		{Name: "nginx:latest", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
		{Name: "nginx:latest", ImageLocations: []types.ImageLocation{{Origin: types.DockerComposeFileOrigin, Path: "docker-compose.yaml"}}},
		{Name: "mysql:5.7", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "values.yaml"}}},
		{Name: "mysql:5.7", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "values.yaml"}}},
	}

	// Call the mergeDuplicates function
	mergedImages := mergeDuplicates(imageModels)

	// Define the expected result with duplicates merged
	expectedResult := []types.ImageModel{
		{Name: "nginx:latest", ImageLocations: []types.ImageLocation{
			{Origin: types.DockerFileOrigin, Path: "Dockerfile"},
			{Origin: types.DockerComposeFileOrigin, Path: "docker-compose.yaml"},
		}},
		{Name: "mysql:5.7", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "values.yaml"}}},
	}

	// Check if the merged images match the expected result
	if !reflect.DeepEqual(mergedImages, expectedResult) {
		t.Errorf("MergeDuplicates result does not match the expected result")
	}
}
