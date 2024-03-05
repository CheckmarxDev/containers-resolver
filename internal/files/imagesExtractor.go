package files

import (
	"github.com/Checkmarx-Containers/containers-resolver/internal/extractors"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
	"log"
)

func MergeImages(images, imagesFromFiles, helmImages []types.ImageModel) []types.ImageModel {
	if len(imagesFromFiles) > 0 {
		images = append(images, imagesFromFiles...)
	}
	if len(helmImages) > 0 {
		images = append(images, helmImages...)
	}
	return removeDuplicates(images)
}

func ExtractImagesFromFiles(files types.FileImages) ([]types.ImageModel, error) {

	dockerfileImages, err := extractors.ExtractImagesFromDockerfiles(files.Dockerfile)
	if err != nil {
		log.Println("Could not extract images from docker files", err)
		return nil, err
	}

	dockerComposeFileImages, err := extractors.ExtractImagesFromDockerComposeFiles(files.DockerCompose)
	if err != nil {
		log.Println("Could not extract images from docker compose files", err)
		return nil, err
	}

	helmImages, err := extractors.ExtractImagesFromHelmFiles(files.Helm)
	if err != nil {
		log.Println("Could not extract images from helm files", err)
		return nil, err
	}

	imagesFromFiles := MergeImages(dockerfileImages, dockerComposeFileImages, helmImages)

	return imagesFromFiles, nil
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
