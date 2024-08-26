package containersResolver

import (
	"github.com/CheckmarxDev/containers-resolver/internal/files"
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	se "github.com/CheckmarxDev/containers-resolver/internal/syftExtractor"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
)

func Resolve(scanPath string, resolutionFolderPath string, images []string, isDebug bool) error {

	resolverLogger := logger.NewLogger(isDebug)

	imagesExtractor := files.ImagesExtractor{
		Logger: resolverLogger,
	}

	syftExtractor := se.SyftExtractor{
		Logger: resolverLogger,
	}
	resolverLogger.Debug("Resolve func parameters: scanPath=%s, resolutionFolderPath=%s, images=%s, isDebug=%t", scanPath, resolutionFolderPath, images, isDebug)

	// 0. validate input
	err := validate(resolutionFolderPath)
	if err != nil {
		resolverLogger.Error("input is not valid. err: %v", err)
		return err
	}

	//1. extract files
	filesWithImages, settingsFiles, outputPath, err := imagesExtractor.ExtractFiles(scanPath)
	if err != nil {
		resolverLogger.Error("Could not extract files. err: %v", err)
		return err
	}

	//2. extract images from files
	imagesToAnalyze, err := imagesExtractor.ExtractAndMergeImagesFromFiles(filesWithImages, toImageModels(images), settingsFiles)
	if err != nil {
		resolverLogger.Error("Could not extract images from files err: %+v", err)
		return err
	}

	//4. get images resolution
	resolutionResult, err := syftExtractor.AnalyzeImages(imagesToAnalyze)
	if err != nil {
		resolverLogger.Error("Could not analyze images. err: %v", err)
		return err
	}

	//5. save to resolution file path
	err = imagesExtractor.SaveObjectToFile(resolutionFolderPath, resolutionResult)
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
	isValidFolderPath, err := files.IsValidFolderPath(resolutionFolderPath)
	if err != nil || isValidFolderPath == false {
		return err
	}
	return nil
}

func cleanup(originalPath string, outputPath string) error {
	if outputPath != "" && outputPath != originalPath {
		err := files.DeleteDirectory(outputPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func toImageModels(images []string) []types.ImageModel {
	imageNames := []types.ImageModel{}

	for _, image := range images {
		imageNames = append(imageNames, types.ImageModel{
			Name: image,
			ImageLocations: []types.ImageLocation{
				{
					Origin: types.UserInput,
					Path:   types.NoFilePath,
				},
			},
		})
	}

	return imageNames
}
