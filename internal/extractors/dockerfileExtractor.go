package extractors

import (
	"bufio"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
	"log"
	"os"
	"regexp"
	"strings"
)

func ExtractImagesFromDockerfiles(filePaths []types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	for _, filePath := range filePaths {
		fileImages, err := extractImagesFromDockerfile(filePath)
		if err != nil {
			return nil, err
		}

		imageNames = append(imageNames, fileImages...)
	}

	return imageNames, nil
}

func extractImagesFromDockerfile(filePath types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	file, err := os.Open(filePath.FullPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Could not close dockerfile:", file.Name())
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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	printFoundImagesInFile(filePath.RelativePath, imageNames)

	return imageNames, nil
}

func printFoundImagesInFile(filePath string, imageNames []types.ImageModel) {
	if len(imageNames) > 0 {
		log.Printf("Successfully found images in file: %s images are: %v\n", filePath, strings.Join(func() []string {
			var result []string
			for _, obj := range imageNames {
				result = append(result, obj.Name)
			}
			return result
		}(), ", "))

	} else {
		log.Printf("Could not find any images in file: %s\n", filePath)
	}
}
