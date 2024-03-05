package files

import (
	"encoding/json"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
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

func ExtractFiles(scanPath string) (types.FileImages, string, error) {

	filesPath, err := extractCompressedPath(scanPath)
	if err != nil {
		log.Printf("Could not extract compressed folder: %s", err.Error())
		return types.FileImages{}, scanPath, err
	}

	var f types.FileImages

	err = filepath.Walk(filesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path matches the Dockerfile pattern
		if dockerfilePattern.MatchString(info.Name()) {
			f.Dockerfile = append(f.Dockerfile, types.FilePath{
				FullPath:     path,
				RelativePath: getRelativePath(filesPath, path),
			})
		}

		// Check if the current path matches the Docker Compose file pattern
		if dockerComposePattern.MatchString(info.Name()) {
			f.DockerCompose = append(f.DockerCompose, types.FilePath{
				FullPath:     path,
				RelativePath: getRelativePath(filesPath, path),
			})
		}

		return nil
	})

	if err != nil {
		log.Printf("Could not extract docker or docker compose files: %s", err.Error())
	}

	helmCharts, err := findHelmCharts(filesPath)
	if err != nil {
		log.Printf("Could not extract helm charts: %s", err.Error())
	}
	if len(helmCharts) > 0 {
		f.Helm = helmCharts
	}
	printFilePaths(f.Dockerfile, "Successfully found dockerfiles")
	printFilePaths(f.DockerCompose, "Successfully found docker compose files")

	return f, filesPath, nil
}

func IsValidFolderPath(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func getContainerResolutionFullPath(folderPath string) (string, error) {
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

func DeleteDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		log.Printf("Could Not Delete directory %s", dirPath)
		return err
	}
	return nil
}

func isHelmChart(directory string) bool {
	chartFilePath := filepath.Join(directory, "Chart.yaml")
	valuesFilePath := filepath.Join(directory, "values.yaml")
	templatesDirPath := filepath.Join(directory, "templates")

	_, errChart := os.Stat(chartFilePath)
	_, errValues := os.Stat(valuesFilePath)
	_, errTemplatesDir := os.Stat(templatesDirPath)

	return errChart == nil && errValues == nil && errTemplatesDir == nil
}

func findHelmCharts(baseDir string) ([]types.HelmChartInfo, error) {
	var helmCharts []types.HelmChartInfo

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && isHelmChart(path) {

			valuesFile := filepath.Join(path, "values.yaml")
			relativeValuesPath, _ := filepath.Rel(baseDir, valuesFile)

			templatesDir := filepath.Join(path, "templates")

			var templateFiles []types.FilePath
			err := filepath.Walk(templatesDir, func(templatePath string, templateInfo os.FileInfo, templateErr error) error {
				if templateErr != nil {
					return templateErr
				}
				if !templateInfo.IsDir() && isYAMLFile(templatePath) {
					relativeTemplatePath, _ := filepath.Rel(baseDir, templatePath)
					templateFiles = append(templateFiles, types.FilePath{
						FullPath:     templatePath,
						RelativePath: relativeTemplatePath,
					})
				}

				return nil
			})

			if err != nil {
				return err
			}

			helmChart := types.HelmChartInfo{
				Directory:     path,
				ValuesFile:    relativeValuesPath,
				TemplateFiles: templateFiles,
			}

			helmCharts = append(helmCharts, helmChart)
		}

		return nil
	})

	return helmCharts, err
}

func isYAMLFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".yml" || ext == ".yaml"
}

func extractCompressedPath(inputPath string) (string, error) {
	if fileInfo, err := os.Stat(inputPath); err == nil && fileInfo.IsDir() {
		return inputPath, nil
	}

	if strings.HasSuffix(inputPath, ".zip") {
		return extractZip(inputPath)
	}

	if strings.HasSuffix(inputPath, ".tar") || strings.HasSuffix(inputPath, ".tar.gz") || strings.HasSuffix(inputPath, ".tgz") {
		return extractTar(inputPath)
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

func printFilePaths(f []types.FilePath, message string) {
	log.Printf("%s. files: %v\n", message, strings.Join(func() []string {
		var result []string
		for _, obj := range f {
			result = append(result, obj.RelativePath)
		}
		return result
	}(), ", "))
}
