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
	return removeDuplicates(images)
}

func removeDuplicates(slice []types.ImageModel) []types.ImageModel {
	seen := make(map[types.ImageModel]bool)
	var result []types.ImageModel

	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}
