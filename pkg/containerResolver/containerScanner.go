package containersResolver

import (
	"github.com/Checkmarx-Containers/extractor-types/types"
	"github.com/Checkmarx-Containers/images-extractor/pkg/imagesExtractor"
	"github.com/Checkmarx-Containers/syft-packages-extractor/pkg/syftPackagesExtractor"
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
)

type ContainersResolver struct {
	imagesExtractor.ImagesExtractorInterface
	syftPackagesExtractor.SyftPackagesExtractorInterface
}

func (cr ContainersResolver) Resolve(scanPath string, resolutionFolderPath string, images []string, isDebug bool) error {

	resolverLogger := logger.NewLogger(isDebug)

	resolverLogger.Debug("Resolve func parameters: scanPath=%s, resolutionFolderPath=%s, images=%s, isDebug=%t", scanPath, resolutionFolderPath, images, isDebug)

	// 0. validate input
	err := validate(resolutionFolderPath)
	if err != nil {
		resolverLogger.Error("input is not valid. err: %v", err)
		return err
	}

	//1. extract files
	filesWithImages, settingsFiles, outputPath, err := cr.ExtractFiles(scanPath)
	if err != nil {
		resolverLogger.Error("Could not extract files. err: %v", err)
		return err
	}

	//2. extract images from files
	imagesToAnalyze, err := cr.ExtractAndMergeImagesFromFiles(filesWithImages, types.ToImageModels(images), settingsFiles)
	if err != nil {
		resolverLogger.Error("Could not extract images from files err: %+v", err)
		return err
	}

	//4. get images resolution
	resolutionResult, err := cr.AnalyzeImages(imagesToAnalyze)
	if err != nil {
		resolverLogger.Error("Could not analyze images. err: %v", err)
		return err
	}

	//5. save to resolution file path
	err = cr.SaveObjectToFile(resolutionFolderPath, resolutionResult)
	if err != nil {
		resolverLogger.Error("Could not save resolution result. err: %v", err)
		return err
	}
	//6. cleanup files generated folder
	err = cleanup(resolutionFolderPath, outputPath)
	if err != nil {
		resolverLogger.Error("Could not cleanup resources. err: %v", err)
		return err
	}
	return nil
}

func validate(resolutionFolderPath string) error {
	isValidFolderPath, err := imagesExtractor.IsValidFolderPath(resolutionFolderPath)
	if err != nil || isValidFolderPath == false {
		return err
	}
	return nil
}

func cleanup(originalPath string, outputPath string) error {
	if outputPath != "" && outputPath != originalPath {
		err := imagesExtractor.DeleteDirectory(outputPath)
		if err != nil {
			return err
		}
	}
	return nil
}
