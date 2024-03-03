package files

import (
	"encoding/json"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/files/zip"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	dockerfilePattern    = regexp.MustCompile(`Dockerfile$`)
	dockerComposePattern = regexp.MustCompile(`docker-compose(\.yml|\.yaml)$`)
)

func ExtractFiles(scanPath string) (FileImages, error) {

	filesPath, err := extractCompressedPath(scanPath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var f FileImages

	err = filepath.Walk(filesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path matches the Dockerfile pattern
		if dockerfilePattern.MatchString(info.Name()) {
			f.Dockerfile = append(f.Dockerfile, FilePath{
				FullPath:     path,
				RelativePath: getRelativePath(filesPath, path),
			})
		}

		// Check if the current path matches the Docker Compose file pattern
		if dockerComposePattern.MatchString(info.Name()) {
			f.DockerCompose = append(f.DockerCompose, FilePath{
				FullPath:     path,
				RelativePath: getRelativePath(filesPath, path),
			})
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	printFilePaths(f.Dockerfile, "Successfully found dockerfiles")
	printFilePaths(f.DockerCompose, "Successfully found docker compose files")

	return f, err
}

func isValidFolderPath(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func getContainerResolutionFullPath(folderPath string) (string, error) {
	if !isValidFolderPath(folderPath) {
		return "", fmt.Errorf("invalid folder path: %s", folderPath)
	}
	return folderPath + "/containers-resolution.json", nil // Hard-coding the containers resolution filename
}

func SaveObjectToFile(folderPath string, obj interface{}) error {
	containerResolutionFullPath, err := getContainerResolutionFullPath(folderPath)
	if err != nil {
		fmt.Println("Error getting container resolution full file path:", err)
		return err
	}
	fmt.Println("containers-resolution.json full path:", containerResolutionFullPath)

	resultBytes, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Error marshaling struct:", err)
		return err
	}

	err = os.WriteFile(containerResolutionFullPath, resultBytes, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}

func extractCompressedPath(inputPath string) (string, error) {
	if fileInfo, err := os.Stat(inputPath); err == nil && fileInfo.IsDir() {
		return inputPath, nil
	}

	if strings.HasSuffix(inputPath, ".zip") {
		return zip.ExtractZip(inputPath)
	}

	if strings.HasSuffix(inputPath, ".tar") || strings.HasSuffix(inputPath, ".tar.gz") || strings.HasSuffix(inputPath, ".tgz") {
		return zip.ExtractTar(inputPath)
	}

	return "", fmt.Errorf("unsupported file type: %s", inputPath)
}

func getRelativePath(baseDir, filePath string) string {
	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return filePath
	}
	return relativePath
}

func printFilePaths(f []FilePath, message string) {
	log.Printf("%s. files: %v\n", message, strings.Join(func() []string {
		var result []string
		for _, obj := range f {
			result = append(result, obj.RelativePath)
		}
		return result
	}(), ", "))
}
