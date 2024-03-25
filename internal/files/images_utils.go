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
	var result []types.ImageModel

	for _, img := range imageModels {
		if _, ok := aggregated[img.Name]; !ok {
			// If the image name is not yet in the result, add it with its locations
			result = append(result, types.ImageModel{Name: img.Name, ImageLocations: img.ImageLocations})
			aggregated[img.Name] = img.ImageLocations
		} else {
			// If the image name is already in the result, merge the locations
			for _, location := range img.ImageLocations {
				found := false
				for _, existingLocation := range aggregated[img.Name] {
					if existingLocation.Origin == location.Origin && existingLocation.Path == location.Path {
						found = true
						break
					}
				}
				if !found {
					// Append only new locations to the existing entry in the result slice
					for i := range result {
						if result[i].Name == img.Name {
							result[i].ImageLocations = append(result[i].ImageLocations, location)
							break
						}
					}
					// Update the map to include the new location
					aggregated[img.Name] = append(aggregated[img.Name], location)
				}
			}
		}
	}

	return result
}
