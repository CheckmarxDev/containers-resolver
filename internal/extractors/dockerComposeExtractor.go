package extractors

import (
	"bufio"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
	"log"
	"os"
	"regexp"
)

func ExtractImagesFromDockerComposeFiles(filePaths []types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	for _, filePath := range filePaths {
		fileImages, err := extractImagesFromDockerComposeFile(filePath)
		if err != nil {
			return nil, err
		}

		imageNames = append(imageNames, fileImages...)
	}

	return imageNames, nil
}

func extractImagesFromDockerComposeFile(filePath types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	file, err := os.Open(filePath.FullPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Could not close docker compose file:", file.Name())
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if match := regexp.MustCompile(`^\s*image:\s*([\w./-]+)(?::([\w.-]+))?`).FindStringSubmatch(line); match != nil {
			imageName := match[1]
			tag := match[2]

			if tag == "" {
				tag = "latest"
			}

			fullImageName := fmt.Sprintf("%s:%s", imageName, tag)

			imageNames = append(imageNames, types.ImageModel{
				Name:   fullImageName,
				Origin: types.DockerComposeFileOrigin,
				Path:   filePath.RelativePath,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	printFoundImagesInFile(filePath.RelativePath, imageNames)

	return imageNames, nil
}
