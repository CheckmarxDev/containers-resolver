package extractors

import (
	"bufio"
	"fmt"
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
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
			logger.Warn("could not extract images from dockerfile %s err: %+v", filePath, err)
		}
		printFoundImagesInFile(logger, filePath.RelativePath, fileImages)
		imageNames = append(imageNames, fileImages...)
	}

	return imageNames, nil
}

func extractImagesFromDockerfile(logger *logger.Logger, filePath types.FilePath) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	aliases := make(map[string]string) // Map to store aliases and their corresponding real image names

	file, err := os.Open(filePath.FullPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Warn("Could not close dockerfile: %s err: %+v", file.Name(), err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line defines an alias
		if match := regexp.MustCompile(`^\s*FROM\s+([\w./-]+(?::[\w.-]+)?)(?:\s+AS\s+(\w+))?`).FindStringSubmatch(line); match != nil {
			imageName := match[1]
			alias := match[2]
			if alias != "" {
				// Check if the alias points to another alias
				realName := resolveAlias(alias, aliases)
				if realName != "" {
					aliases[alias] = realName // Store the alias and its corresponding real image name
				} else {
					aliases[alias] = imageName // Store the alias and its corresponding real image name
				}
			}
		}

		// Check if the line contains an image reference
		if match := regexp.MustCompile(`\bFROM\s+([\w./-]+)(?::([\w.-]+))?\b`).FindStringSubmatch(line); match != nil {
			imageName := match[1]
			tag := match[2]

			if tag == "" {
				tag = "latest"
			}

			fullImageName := fmt.Sprintf("%s:%s", imageName, tag)

			if realName, ok := aliases[imageName]; ok {
				if realName != imageName {
					continue
				}
			}

			imageNames = append(imageNames, types.ImageModel{
				Name: fullImageName,
				ImageLocations: []types.ImageLocation{
					{
						Origin: types.DockerFileOrigin,
						Path:   filePath.RelativePath,
					},
				},
			})
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return imageNames, nil
}

func resolveAlias(alias string, aliases map[string]string) string {
	realName, ok := aliases[alias]
	if !ok {
		return "" // Alias not found
	}

	// Check if the real name is also an alias, resolve recursively
	if resolvedRealName, ok := aliases[realName]; ok {
		return resolveAlias(resolvedRealName, aliases)
	}

	return realName
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
