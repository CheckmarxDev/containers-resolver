package files

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestExtractAndMergeImagesFromFiles(t *testing.T) {
	// Initialize logger
	l := logger.NewLogger(false)
	extractor := &ImagesExtractor{Logger: l}

	// Define test scenarios
	scenarios := []struct {
		Name             string
		Files            types.FileImages
		UserInput        []types.ImageModel
		ExpectedImages   []types.ImageModel
		ExpectedErrorMsg string
	}{
		{
			Name: "DifferentImagesFromDifferentSources",
			Files: types.FileImages{
				Dockerfile: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile", RelativePath: "Dockerfile"},
				},
				DockerCompose: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose.yaml", RelativePath: "docker-compose1.yml"},
				},
				Helm: []types.HelmChartInfo{
					{
						Directory:  "../../test_files/imageExtraction/helm/",
						ValuesFile: "../../test_files/imageExtraction/helm/values.yaml",
						TemplateFiles: []types.FilePath{{FullPath: "../test_files/imageExtraction/helm/templates/containers-worker.yaml", RelativePath: "templates/containers-worker.yaml"},
							{FullPath: "../test_files/imageExtraction/helm/templates/image-insights.yaml", RelativePath: "templates/image-insights.yaml"}},
					},
				},
			},
			UserInput: []types.ImageModel{
				{Name: "debian:11", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}},
			},
			ExpectedImages: []types.ImageModel{
				{Name: "debian:11", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}},
				{Name: "mcr.microsoft.com/dotnet/sdk:6.0", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
				{Name: "mcr.microsoft.com/dotnet/aspnet:6.0", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
				{Name: "buildimage:latest", ImageLocations: []types.ImageLocation{{Origin: types.DockerComposeFileOrigin, Path: "docker-compose1.yml"}}},
				{Name: "checkmarx.jfrog.io/ast-docker/containers-worker:b201b1f", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "containers/templates/containers-worker.yaml"}}},
				{Name: "checkmarx.jfrog.io/ast-docker/image-insights:f4b507b", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "containers/templates/image-insights.yaml"}}},
			},
		},
		{
			Name: "SameImagesFromDifferentSources",
			Files: types.FileImages{
				Dockerfile: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile", RelativePath: "Dockerfile"},
				},
				DockerCompose: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose-4.yaml", RelativePath: "docker-compose-4.yml"},
				},
			},
			UserInput: []types.ImageModel{
				{Name: "mcr.microsoft.com/dotnet/sdk:6.0", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}},
			},
			ExpectedImages: []types.ImageModel{
				{Name: "mcr.microsoft.com/dotnet/sdk:6.0", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath},
					{Origin: types.DockerFileOrigin, Path: "Dockerfile"}, {Origin: types.DockerComposeFileOrigin, Path: "docker-compose-4.yml"}}},
				{Name: "mcr.microsoft.com/dotnet/aspnet:6.0", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}}},
		},
		{
			Name: "OnlyDockerfileFound",
			Files: types.FileImages{
				Dockerfile: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile", RelativePath: "Dockerfile"},
				},
			},
			UserInput: []types.ImageModel{
				{Name: "debian:11", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}},
			},
			ExpectedImages: []types.ImageModel{
				{Name: "mcr.microsoft.com/dotnet/sdk:6.0", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
				{Name: "mcr.microsoft.com/dotnet/aspnet:6.0", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
				{Name: "debian:11", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}}},
		},
		{
			Name: "OnlyDockerComposeFound",
			Files: types.FileImages{
				DockerCompose: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose-4.yaml", RelativePath: "docker-compose-4.yml"},
				},
			},
			UserInput: []types.ImageModel{
				{Name: "debian:11", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}},
			},
			ExpectedImages: []types.ImageModel{
				{Name: "mcr.microsoft.com/dotnet/sdk:6.0", ImageLocations: []types.ImageLocation{{Origin: types.DockerComposeFileOrigin, Path: "docker-compose-4.yml"}}},
				{Name: "debian:11", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}}},
		},
		{
			Name: "OnlyHelmChartFound",
			Files: types.FileImages{
				Helm: []types.HelmChartInfo{
					{
						Directory:  "../../test_files/imageExtraction/helm/",
						ValuesFile: "../../test_files/imageExtraction/helm/values.yaml",
						TemplateFiles: []types.FilePath{{FullPath: "../test_files/imageExtraction/helm/templates/containers-worker.yaml", RelativePath: "templates/containers-worker.yaml"},
							{FullPath: "../test_files/imageExtraction/helm/templates/image-insights.yaml", RelativePath: "templates/image-insights.yaml"}},
					},
				},
			},
			UserInput: []types.ImageModel{
				{Name: "debian:11", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}},
			},
			ExpectedImages: []types.ImageModel{
				{Name: "debian:11", ImageLocations: []types.ImageLocation{{Origin: types.UserInput, Path: types.NoFilePath}}},
				{Name: "checkmarx.jfrog.io/ast-docker/containers-worker:b201b1f", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "containers/templates/containers-worker.yaml"}}},
				{Name: "checkmarx.jfrog.io/ast-docker/image-insights:f4b507b", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "containers/templates/image-insights.yaml"}}},
			},
		},
		{
			Name: "AllTypesOfFilesWithNoExistingImages",
			Files: types.FileImages{
				Dockerfile: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile", RelativePath: "Dockerfile"},
				},
				DockerCompose: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose.yaml", RelativePath: "docker-compose1.yml"},
				},
				Helm: []types.HelmChartInfo{
					{
						Directory:  "../../test_files/imageExtraction/helm/",
						ValuesFile: "../../test_files/imageExtraction/helm/values.yaml",
						TemplateFiles: []types.FilePath{{FullPath: "../test_files/imageExtraction/helm/templates/containers-worker.yaml", RelativePath: "templates/containers-worker.yaml"},
							{FullPath: "../test_files/imageExtraction/helm/templates/image-insights.yaml", RelativePath: "templates/image-insights.yaml"}},
					},
				},
			},
			ExpectedImages: []types.ImageModel{
				{Name: "mcr.microsoft.com/dotnet/sdk:6.0", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
				{Name: "mcr.microsoft.com/dotnet/aspnet:6.0", ImageLocations: []types.ImageLocation{{Origin: types.DockerFileOrigin, Path: "Dockerfile"}}},
				{Name: "buildimage:latest", ImageLocations: []types.ImageLocation{{Origin: types.DockerComposeFileOrigin, Path: "docker-compose1.yml"}}},
				{Name: "checkmarx.jfrog.io/ast-docker/containers-worker:b201b1f", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "containers/templates/containers-worker.yaml"}}},
				{Name: "checkmarx.jfrog.io/ast-docker/image-insights:f4b507b", ImageLocations: []types.ImageLocation{{Origin: types.HelmFileOrigin, Path: "containers/templates/image-insights.yaml"}}},
			},
		},
	}

	// Run test scenarios
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// Run the function
			result, err := extractor.ExtractAndMergeImagesFromFiles(scenario.Files, scenario.UserInput, nil)

			// Check for errors
			if scenario.ExpectedErrorMsg != "" {
				if err == nil || err.Error() != scenario.ExpectedErrorMsg {
					t.Errorf("Expected error message '%s' but got '%v'", scenario.ExpectedErrorMsg, err)
				}
			} else {
				// Check for expected images
				expectedImageMap := make(map[string][]types.ImageLocation)
				for _, img := range scenario.ExpectedImages {
					expectedImageMap[img.Name] = img.ImageLocations
				}

				if len(result) != len(scenario.ExpectedImages) {
					t.Errorf("Expected %d images but got %d", len(scenario.ExpectedImages), len(result))
				}

				for _, img := range result {
					expectedLocations, exists := expectedImageMap[img.Name]
					if !exists {
						t.Errorf("Unexpected image found: %s", img.Name)
						continue
					}

					if !reflect.DeepEqual(img.ImageLocations, expectedLocations) {
						t.Errorf("Image locations mismatch for image '%s'", img.Name)
					}

					// Remove processed image from map
					delete(expectedImageMap, img.Name)
				}

				// Check for any remaining expected images
				if len(expectedImageMap) > 0 {
					for name := range expectedImageMap {
						t.Errorf("Expected image not found: %s", name)
					}
				}
			}
		})
	}
}

func TestExtractFiles(t *testing.T) {
	// Initialize logger
	l := logger.NewLogger(false)
	extractor := &ImagesExtractor{Logger: l}

	// Define test scenarios
	scenarios := []struct {
		Name              string
		InputPath         string
		ExpectedFiles     types.FileImages
		ExpectedErrString string
	}{
		{
			Name:      "FolderInput",
			InputPath: "../../test_files/imageExtraction",
			ExpectedFiles: types.FileImages{
				Dockerfile: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile", RelativePath: "dockerfiles/Dockerfile"},
					{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile-2", RelativePath: "dockerfiles/Dockerfile-2"},
					{FullPath: "../../test_files/imageExtraction/dockerfiles/Dockerfile-3", RelativePath: "dockerfiles/Dockerfile-3"},
				},
				DockerCompose: []types.FilePath{
					{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose.yaml", RelativePath: "dockerCompose/docker-compose.yaml"},
					{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose-2.yaml", RelativePath: "dockerCompose/docker-compose-2.yaml"},
					{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose-3.yaml", RelativePath: "dockerCompose/docker-compose-3.yaml"},
					{FullPath: "../../test_files/imageExtraction/dockerCompose/docker-compose-4.yaml", RelativePath: "dockerCompose/docker-compose-4.yaml"},
				},
				Helm: []types.HelmChartInfo{
					{
						Directory:  "../../test_files/imageExtraction/helm",
						ValuesFile: "helm/values.yaml",
						TemplateFiles: []types.FilePath{
							{FullPath: "../../test_files/imageExtraction/helm/templates/containers-worker.yaml", RelativePath: "helm/templates/containers-worker.yaml"},
							{FullPath: "../../test_files/imageExtraction/helm/templates/image-insights.yaml", RelativePath: "helm/templates/image-insights.yaml"},
						},
					},
				},
			},
			ExpectedErrString: "",
		},
		{
			Name:      "TarInput",
			InputPath: "../../test_files/withDockerInTar.tar.gz",
			ExpectedFiles: types.FileImages{
				Dockerfile: []types.FilePath{
					{FullPath: "../../test_files/extracted_tar/withDockerInTar/Dockerfile", RelativePath: "withDockerInTar/Dockerfile"},
					{FullPath: "../../test_files/extracted_tar/withDockerInTar/integrationTests/Dockerfile", RelativePath: "withDockerInTar/integrationTests/Dockerfile"},
				},
				DockerCompose: []types.FilePath{
					{FullPath: "../../test_files/extracted_tar/withDockerInTar/docker-compose.yaml", RelativePath: "withDockerInTar/docker-compose.yaml"},
				},
			},
			ExpectedErrString: "",
		},
		{
			Name:      "ZipInput",
			InputPath: "../../test_files/withDockerInZip.zip",
			ExpectedFiles: types.FileImages{
				Dockerfile: []types.FilePath{
					{FullPath: "../../test_files/extracted_zip/Dockerfile", RelativePath: "Dockerfile"},
					{FullPath: "../../test_files/extracted_zip/integrationTests/Dockerfile", RelativePath: "integrationTests/Dockerfile"},
				},
				DockerCompose: []types.FilePath{
					{FullPath: "../../test_files/extracted_zip/docker-compose.yaml", RelativePath: "docker-compose.yaml"},
				},
			},
			ExpectedErrString: "",
		},
	}

	// Run test scenarios
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// Run the function
			files, _, _, err := extractor.ExtractFiles(scenario.InputPath)

			// Check for errors
			if scenario.ExpectedErrString != "" {
				if err == nil || !strings.Contains(err.Error(), scenario.ExpectedErrString) {
					t.Errorf("Expected error containing '%s' but got '%v'", scenario.ExpectedErrString, err)
				}
			} else {
				if !CompareDockerfiles(files.Dockerfile, scenario.ExpectedFiles.Dockerfile) {
					t.Errorf("Extracted Dockerfiles mismatch for scenario '%s'", scenario.Name)
				}
				if !CompareDockerCompose(files.DockerCompose, scenario.ExpectedFiles.DockerCompose) {
					t.Errorf("Extracted Docker Compose files mismatch for scenario '%s'", scenario.Name)
				}
				if !CompareHelm(files.Helm, scenario.ExpectedFiles.Helm) {
					t.Errorf("Extracted Helm charts mismatch for scenario '%s'", scenario.Name)
				}
			}
		})
	}
}

func CompareDockerfiles(a, b []types.FilePath) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].RelativePath < a[j].RelativePath
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].RelativePath < b[j].RelativePath
	})
	for i := range a {
		if a[i].FullPath != b[i].FullPath {
			return false
		}
		if a[i].RelativePath != b[i].RelativePath {
			return false
		}
	}
	return true
}

// CompareDockerCompose compares two slices of FilePath.
func CompareDockerCompose(a, b []types.FilePath) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].RelativePath < a[j].RelativePath
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].RelativePath < b[j].RelativePath
	})
	for i := range a {
		if a[i].FullPath != b[i].FullPath {
			return false
		}
		if a[i].RelativePath != b[i].RelativePath {
			return false
		}
	}
	return true
}

// CompareHelm compares two slices of HelmChartInfo.
func CompareHelm(a, b []types.HelmChartInfo) bool {
	if len(a) != len(b) {
		return false
	}

	// Sort slices by directory to ensure consistent comparison
	sort.Slice(a, func(i, j int) bool {
		return a[i].Directory < a[j].Directory
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].Directory < b[j].Directory
	})

	// Iterate over each HelmChartInfo struct
	for i := range a {
		// Compare Directory and ValuesFile
		if a[i].Directory != b[i].Directory || a[i].ValuesFile != b[i].ValuesFile {
			return false
		}

		// Compare TemplateFiles slices
		if len(a[i].TemplateFiles) != len(b[i].TemplateFiles) {
			return false
		}

		// Sort TemplateFiles slices by RelativePath for consistent comparison
		sort.Slice(a[i].TemplateFiles, func(j, k int) bool {
			return a[i].TemplateFiles[j].RelativePath < a[i].TemplateFiles[k].RelativePath
		})
		sort.Slice(b[i].TemplateFiles, func(j, k int) bool {
			return b[i].TemplateFiles[j].RelativePath < b[i].TemplateFiles[k].RelativePath
		})

		// Iterate over each FilePath struct in TemplateFiles slice
		for j := range a[i].TemplateFiles {
			// Compare FullPath and RelativePath of each FilePath struct
			if a[i].TemplateFiles[j].FullPath != b[i].TemplateFiles[j].FullPath || a[i].TemplateFiles[j].RelativePath != b[i].TemplateFiles[j].RelativePath {
				return false
			}
		}
	}
	return true
}
