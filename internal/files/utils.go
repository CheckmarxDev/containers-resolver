package files

import (
	"fmt"
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	dockerfilePattern    = regexp.MustCompile(`Dockerfile$`)
	dockerComposePattern = regexp.MustCompile(`docker-compose(\.yml|\.yaml)$`)
)

func IsValidFolderPath(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func DeleteDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	return nil
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
			err = filepath.Walk(templatesDir, func(templatePath string, templateInfo os.FileInfo, templateErr error) error {
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

func extractCompressedPath(l *logger.Logger, inputPath string) (string, error) {
	if fileInfo, err := os.Stat(inputPath); err == nil && fileInfo.IsDir() {
		return inputPath, nil
	}

	if strings.HasSuffix(inputPath, ".zip") {
		return extractZip(l, inputPath)
	}

	if strings.HasSuffix(inputPath, ".tar") || strings.HasSuffix(inputPath, ".tar.gz") || strings.HasSuffix(inputPath, ".tgz") {
		return extractTar(l, inputPath)
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

func isYAMLFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".yml" || ext == ".yaml"
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

func getContainerResolutionFullPath(folderPath string) (string, error) {
	return folderPath + "/containers-resolution.json", nil // Hard-coding the containers resolution filename
}

func printFilePaths(logger *logger.Logger, f []types.FilePath, message string) {
	logger.Debug("%s. files: %v\n", message, strings.Join(func() []string {
		var result []string
		for _, obj := range f {
			result = append(result, obj.RelativePath)
		}
		return result
	}(), ", "))
}
