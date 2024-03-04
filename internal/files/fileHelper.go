package files

import (
	"encoding/json"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/zip"
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

func SaveObjectToFile(filePath string, obj interface{}) error {

	resultBytes, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Error marshaling struct:", err)
		return err
	}

	err = os.WriteFile(filePath, resultBytes, 0644)
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
