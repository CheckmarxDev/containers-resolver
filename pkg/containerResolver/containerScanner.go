package containersResolver

import (
	"github.com/Checkmarx-Containers/containers-resolver/internal/files"
	"github.com/Checkmarx-Containers/containers-resolver/internal/syftExtractor"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
	"log"
)

func Resolve(scanPath string, resolutionFolderPath string, images []string, isDebug bool) error {
	log.Printf("Resolve func parameters: scanPath=%s, resolutionFolderPath=%s, images=%s, isDebug=%t", scanPath, resolutionFolderPath, images, isDebug)
	err := validate(resolutionFolderPath)

	//1. extract files
	filesWithImages, outputPath, err := files.ExtractFiles(scanPath)
	if err != nil {
		log.Fatal("Could not extract files", err)
		return err
	}

	//2. extract images from files
	imagesToAnalyze, err := files.ExtractAndMergeImagesFromFiles(filesWithImages, toImageModels(images))
	if err != nil {
		log.Fatal("Could not extract images from files", err)
		return err
	}

	//4. get images resolution
	resolutionResult, err := syftExtractor.AnalyzeImages(imagesToAnalyze)
	if err != nil {
		log.Fatal("Could not analyze images", err)
		return err
	}

	//5. save to resolution file path
	err = files.SaveObjectToFile(resolutionFolderPath, resolutionResult)
	if err != nil {
		log.Fatal("Could not save resolution result", err)
		return err
	}
	//6. cleanup files generated folder
	err = cleanup(resolutionFolderPath, outputPath)
	if err != nil {
		log.Fatal("Could not cleanup resources", err)
		return err
	}
	return nil
}

func validate(resolutionFolderPath string) error {
	isValidFolderPath, err := files.IsValidFolderPath(resolutionFolderPath)
	if err != nil || isValidFolderPath == false {
		log.Fatal("resolutionFolderPath is not a valid path.", err)
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
	var imageNames []types.ImageModel

	for _, image := range images {
		imageNames = append(imageNames, types.ImageModel{
			Name:   image,
			Origin: types.UserInput,
			Path:   types.NoFilePath,
		})
	}

	return imageNames
}
