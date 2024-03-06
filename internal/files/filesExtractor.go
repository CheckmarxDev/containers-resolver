package files

import (
	"encoding/json"
	"github.com/CheckmarxDev/containers-resolver/internal/extractors"
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"os"
	"path/filepath"
)

type ImagesExtractor struct {
	*logger.Logger
}

func (fe *ImagesExtractor) ExtractAndMergeImagesFromFiles(files types.FileImages, images []types.ImageModel) ([]types.ImageModel, error) {

	dockerfileImages, err := extractors.ExtractImagesFromDockerfiles(fe.Logger, files.Dockerfile)
	if err != nil {
		fe.Logger.Error("Could not extract images from docker files", err)
		return nil, err
	}

	dockerComposeFileImages, err := extractors.ExtractImagesFromDockerComposeFiles(fe.Logger, files.DockerCompose)
	if err != nil {
		fe.Logger.Error("Could not extract images from docker compose files", err)
		return nil, err
	}

	helmImages, err := extractors.ExtractImagesFromHelmFiles(fe.Logger, files.Helm)
	if err != nil {
		fe.Logger.Error("Could not extract images from helm files", err)
		return nil, err
	}

	imagesFromFiles := mergeImages(images, dockerfileImages, dockerComposeFileImages, helmImages)

	return imagesFromFiles, nil
}

func (fe *ImagesExtractor) ExtractFiles(scanPath string) (types.FileImages, string, error) {

	filesPath, err := extractCompressedPath(fe.Logger, scanPath)
	if err != nil {
		fe.Logger.Error("Could not extract compressed folder: %s", err)
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
		fe.Logger.Warn("Could not extract docker or docker compose files: %s", err.Error())
	}

	helmCharts, err := findHelmCharts(filesPath)
	if err != nil {
		fe.Logger.Warn("Could not extract helm charts: %s", err.Error())
	}
	if len(helmCharts) > 0 {
		f.Helm = helmCharts
	}
	printFilePaths(fe.Logger, f.Dockerfile, "Successfully found dockerfiles")
	printFilePaths(fe.Logger, f.DockerCompose, "Successfully found docker compose files")

	return f, filesPath, nil
}

func (fe *ImagesExtractor) SaveObjectToFile(folderPath string, obj interface{}) error {
	containerResolutionFullPath, err := getContainerResolutionFullPath(folderPath)
	if err != nil {
		fe.Logger.Error("Error getting container resolution full file path:", err)
		return err
	}
	fe.Logger.Debug("containers-resolution.json full path is: %s", containerResolutionFullPath)

	resultBytes, err := json.Marshal(obj)
	if err != nil {
		fe.Logger.Error("Error marshaling struct:", err)
		return err
	}

	err = os.WriteFile(containerResolutionFullPath, resultBytes, 0644)
	if err != nil {
		fe.Logger.Error("Error writing file:", err)
		return err
	}
	return nil
}
