package extractors

import (
	"bufio"
	"fmt"
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ExtractImagesFromDockerfiles(logger *logger.Logger, filePaths []types.FilePath, envFiles map[string]map[string]string) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel

	for _, filePath := range filePaths {
		logger.Debug("going to extract images from dockerfile %s", filePath)

		fileImages, err := extractImagesFromDockerfile(logger, filePath, envFiles)
		if err != nil {
			logger.Warn("could not extract images from dockerfile %s err: %+v", filePath, err)
		}
		printFoundImagesInFile(logger, filePath.RelativePath, fileImages)
		imageNames = append(imageNames, fileImages...)
	}

	return imageNames, nil
}

func extractImagesFromDockerfile(logger *logger.Logger, filePath types.FilePath, envFiles map[string]map[string]string) ([]types.ImageModel, error) {
	var imageNames []types.ImageModel
	aliases := make(map[string]string)
	argsAndEnv := make(map[string]string)
	mergedEnvVars := resolveEnvVariables(filePath.FullPath, envFiles)

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

		// Parse ARG and ENV lines within the Dockerfile
		if match := regexp.MustCompile(`^\s*(ARG|ENV)\s+(\w+)=([^\s]+)`).FindStringSubmatch(line); match != nil {
			varName := match[2]
			varValue := match[3]
			argsAndEnv[varName] = varValue
		}

		// Replace placeholders with values from mergedEnvVars and argsAndEnv
		line = replacePlaceholders(line, mergedEnvVars, argsAndEnv)

		// Parse FROM instructions
		if match := regexp.MustCompile(`^\s*FROM\s+(?:--platform=[^\s]+\s+)?([\w./-]+(?::[\w.-]+)?)(?:\s+AS\s+(\w+))?`).FindStringSubmatch(line); match != nil {
			imageName := match[1]
			alias := match[2]

			if imageName == "scratch" {
				continue
			}

			if alias != "" {
				realName := resolveAlias(alias, aliases)
				if realName != "" {
					aliases[alias] = realName
				} else {
					aliases[alias] = imageName
				}
			}
		}

		if match := regexp.MustCompile(`\bFROM\s+(?:--platform=[^\s]+\s+)?([\w./-]+)(?::([\w.-]+))?\b`).FindStringSubmatch(line); match != nil {
			imageName := match[1]
			tag := match[2]

			if imageName == "scratch" {
				continue
			}
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

func replacePlaceholders(line string, envVars, argsAndEnv map[string]string) string {
	// Replace ${PLACE} style placeholders first
	for varName, varValue := range envVars {
		placeholderWithBraces := fmt.Sprintf("${%s}", varName)
		line = strings.ReplaceAll(line, placeholderWithBraces, varValue)

		placeholderWithoutBraces := fmt.Sprintf("$%s", varName)
		line = strings.ReplaceAll(line, placeholderWithoutBraces, varValue)
	}

	// Replace ${PLACE} and $PLACE style placeholders for argsAndEnv as well
	for varName, varValue := range argsAndEnv {
		placeholderWithBraces := fmt.Sprintf("${%s}", varName)
		line = strings.ReplaceAll(line, placeholderWithBraces, varValue)

		placeholderWithoutBraces := fmt.Sprintf("$%s", varName)
		line = strings.ReplaceAll(line, placeholderWithoutBraces, varValue)
	}

	return line
}

func resolveEnvVariables(dockerfilePath string, envFiles map[string]map[string]string) map[string]string {
	resolvedVars := make(map[string]string)

	// Iterate over the hierarchy and merge environment variables
	dirs := getDirsForHierarchy(dockerfilePath)
	for _, dir := range dirs {
		if envVars, ok := envFiles[dir]; ok {
			for k, v := range envVars {
				if _, exists := resolvedVars[k]; !exists {
					resolvedVars[k] = v
				}
			}
		}

		//if envVars, ok := envFiles[filepath.Join(dir, ".env")]; ok {
		//	for k, v := range envVars {
		//		if _, exists := resolvedVars[k]; !exists {
		//			resolvedVars[k] = v
		//		}
		//	}
		//}
	}

	return resolvedVars
}

func getDirsForHierarchy(dockerfilePath string) []string {
	var dirs []string

	dir := filepath.Dir(dockerfilePath)
	for dir != "" && dir != "." && dir != "/" {
		dirs = append(dirs, dir)
		dir = filepath.Dir(dir)
	}

	return dirs
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
