package syftExtractor

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
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
		se.Info("successfully analyzed image: %s, found %d packages", imageModel.Name, len(tmpResolution.ContainerPackages))
	}

	if containerResolution == nil || len(containerResolution) < 1 {
		return []*ContainerResolution{}, nil
	}

	return containerResolution, nil
}
