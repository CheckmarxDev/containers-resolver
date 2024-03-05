package extractors

import (
	"bufio"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/logger"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
	"os"
	"regexp"
)

func ExtractImagesFromDockerComposeFiles(logger *logger.Logger, filePaths []types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	for _, filePath := range filePaths {
		logger.Debug("going to extract images from docker compose file %s", filePath)
		fileImages, err := extractImagesFromDockerComposeFile(logger, filePath)
		if err != nil {
			logger.Warn("could not extract images from docker compose file %s", filePath, err)
		}
		printFoundImagesInFile(logger, filePath.RelativePath, fileImages)
		imageNames = append(imageNames, fileImages...)
	}

	return imageNames, nil
}

func extractImagesFromDockerComposeFile(l *logger.Logger, filePath types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	file, err := os.Open(filePath.FullPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			l.Warn("Could not close docker compose file:", file.Name())
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

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return imageNames, nil
}
