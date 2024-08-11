package syftExtractor

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"strings"
)

type SyftExtractor struct {
	*logger.Logger
}

func (se *SyftExtractor) AnalyzeImages(images []types.ImageModel) ([]*ContainerResolution, error) {
	if images == nil {
		return []*ContainerResolution{}, nil
	}

	var containerResolution []*ContainerResolution

	for _, imageModel := range images {
		se.Debug("going to analyze image using syft. image: %s", imageModel.Name)

		tmpResolution, err := analyzeImage(se.Logger, imageModel)
		if err != nil {
			se.Error("Could not analyze image: %s. err: %v", imageModel.Name, err)
			continue
		}
		containerResolution = append(containerResolution, tmpResolution)
		se.Info("successfully analyzed image: %s, found %d packages. image paths: %s", imageModel.Name,
			len(tmpResolution.ContainerPackages), getPaths(imageModel.ImageLocations))

	}

	if containerResolution == nil || len(containerResolution) < 1 {
		return []*ContainerResolution{}, nil
	}

	return containerResolution, nil
}

func getPaths(locations []types.ImageLocation) string {
	var paths []string
	for _, location := range locations {
		paths = append(paths, location.Path)
	}
	return strings.Join(paths, ",")
}
