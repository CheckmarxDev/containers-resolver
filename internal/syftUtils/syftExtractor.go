package syftUtils

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
)

type SyftExtractor struct {
	*logger.Logger
}

func (se *SyftExtractor) AnalyzeImages(images []types.ImageModel) (*ContainerResolution, error) {

	containerResolution := &ContainerResolution{
		ContainerImages:   []ContainerImage{},
		ContainerPackages: []ContainerPackage{},
	}

	for _, imageModel := range images {
		se.Debug("going to analyze image using syft. image: %s", imageModel.Name)

		tmpResolution, err := analyzeImage(se.Logger, imageModel)
		if err != nil {
			se.Error("Could not analyze image: %s err: %+v", imageModel.Name, err)
			continue
		}

		containerResolution.ContainerImages = append(containerResolution.ContainerImages, tmpResolution.ContainerImages...)
		containerResolution.ContainerPackages = append(containerResolution.ContainerPackages, tmpResolution.ContainerPackages...)
		se.Info("successfully analyzed image: %s", imageModel.Name)
	}
	return containerResolution, nil
}
