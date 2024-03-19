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
		if locations, ok := aggregated[img.Name]; ok {
			for _, location := range img.ImageLocations {
				found := false
				for _, existingLocation := range locations {
					if existingLocation.Origin == location.Origin && existingLocation.Path == location.Path {
						found = true
						break
					}
				}
				if !found {
					aggregated[img.Name] = append(aggregated[img.Name], location)
				}
			}
		} else {
			aggregated[img.Name] = img.ImageLocations
		}
	}

	var result []types.ImageModel
	for name, locations := range aggregated {
		result = append(result, types.ImageModel{Name: name, ImageLocations: locations})
	}

	return result
}
