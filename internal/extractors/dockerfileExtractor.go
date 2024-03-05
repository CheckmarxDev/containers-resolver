package extractors

import (
	"bufio"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/logger"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
	"os"
	"regexp"
	"strings"
)

func ExtractImagesFromDockerfiles(logger *logger.Logger, filePaths []types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	for _, filePath := range filePaths {
		logger.Debug("going to extract images from dockerfile %s", filePath)

		fileImages, err := extractImagesFromDockerfile(logger, filePath)
		if err != nil {
			logger.Warn("could not extract images from dockerfile %s", filePath, err)
		}
		printFoundImagesInFile(logger, filePath.RelativePath, fileImages)
		imageNames = append(imageNames, fileImages...)
	}

	return imageNames, nil
}

func extractImagesFromDockerfile(l *logger.Logger, filePath types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	file, err := os.Open(filePath.FullPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			l.Warn("Could not close dockerfile:", file.Name())
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if match := regexp.MustCompile(`\bFROM\s+([\w./-]+)(?::([\w.-]+))?`).FindStringSubmatch(line); match != nil {
			imageName := match[1]
			tag := match[2]

			if tag == "" {
				tag = "latest"
			}

			fullImageName := fmt.Sprintf("%s:%s", imageName, tag)
			imageNames = append(imageNames, types.ImageModel{
				Name:   fullImageName,
				Origin: types.DockerFileOrigin,
				Path:   filePath.RelativePath,
			})
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return imageNames, nil
}

func printFoundImagesInFile(l *logger.Logger, filePath string, imageNames []types.ImageModel) {
	if len(imageNames) > 0 {
		l.Debug("Successfully found images in file: %s images are: %v\n", filePath, strings.Join(func() []string {
			var result []string
			for _, obj := range imageNames {
				result = append(result, obj.Name)
			}
			return result
		}(), ", "))

	} else {
		l.Debug("Could not find any images in file: %s\n", filePath)
	}
}
