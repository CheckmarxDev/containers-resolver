package files

import (
	"github.com/CheckmarxDev/containers-resolver/internal/types"
)

func mergeImages(images, imagesFromDockerFiles, imagesFromDockerComposeFiles, helmImages []types.ImageModel) []types.ImageModel {
	if len(imagesFromDockerFiles) > 0 {
		images = append(images, imagesFromDockerFiles...)
	}
	if len(imagesFromDockerComposeFiles) > 0 {
		images = append(images, imagesFromDockerComposeFiles...)
	}
	if len(helmImages) > 0 {
		images = append(images, helmImages...)
	}
	return mergeDuplicates(images)
}

func mergeDuplicates(imageModels []types.ImageModel) []types.ImageModel {
	aggregated := make(map[string][]types.ImageLocation)

	for _, img := range imageModels {
		aggregated[img.Name] = append(aggregated[img.Name], img.ImageLocations...)
	}

	// Create the final result by constructing ImageModel objects with aggregated ImageLocations
	var result []types.ImageModel
	for name, locations := range aggregated {
		result = append(result, types.ImageModel{Name: name, ImageLocations: locations})
	}

	return result
}
