package files

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func MergeImages(images, imagesFromFiles []ImageModel) []ImageModel {
	if len(imagesFromFiles) > 0 {
		images = append(images, imagesFromFiles...)

	}
	return removeDuplicates(images)
}

func ExtractImagesFromFiles(files FileImages) ([]ImageModel, error) {

	dockerfileImages, err := extractImagesFromDockerfiles(files.Dockerfile)
	if err != nil {
		log.Println("Could not extract images from docker files", err)
		return nil, err
	}

	dockerComposeFileImages, err := extractImagesFromDockerComposeFiles(files.DockerCompose)
	if err != nil {
		log.Println("Could not extract images from docker compose files", err)
		return nil, err
	}

	imagesFromFiles := MergeImages(dockerfileImages, dockerComposeFileImages)

	return removeDuplicates(imagesFromFiles), nil
}

func extractImagesFromDockerComposeFiles(filePaths []FilePath) ([]ImageModel, error) {
	var imageNames []ImageModel

	for _, filePath := range filePaths {
		fileImages, err := extractImagesFromDockerComposeFile(filePath)
		if err != nil {
			return nil, err
		}

		imageNames = append(imageNames, fileImages...)
	}

	return imageNames, nil
}

func extractImagesFromDockerfiles(filePaths []FilePath) ([]ImageModel, error) {
	var imageNames []ImageModel

	for _, filePath := range filePaths {
		fileImages, err := extractImagesFromDockerfile(filePath)
		if err != nil {
			return nil, err
		}

		imageNames = append(imageNames, fileImages...)
	}

	return imageNames, nil
}

func extractImagesFromDockerfile(filePath FilePath) ([]ImageModel, error) {
	var imageNames []ImageModel

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
			imageNames = append(imageNames, ImageModel{
				Name:   fullImageName,
				Origin: DockerFileOrigin,
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

func extractImagesFromDockerComposeFile(filePath FilePath) ([]ImageModel, error) {
	var imageNames []ImageModel

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

			imageNames = append(imageNames, ImageModel{
				Name:   fullImageName,
				Origin: DockerComposeFileOrigin,
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

func printFoundImagesInFile(filePath string, imageNames []ImageModel) {
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

func removeDuplicates(slice []ImageModel) []ImageModel {
	seen := make(map[ImageModel]bool)
	var result []ImageModel

	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}
