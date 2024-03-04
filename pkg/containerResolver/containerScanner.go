package containersResolver

import (
	"github.com/Checkmarx-Containers/containers-resolver/internal/files"
	"github.com/Checkmarx-Containers/containers-resolver/internal/syftExtractor"
	"log"
)

func Resolve(scanPath string, resolutionFilePath string, images []string, isDebug bool) {
	// Print function params
	log.Printf("Resolve func parameters: scanPath=%s, resolutionFilePath=%s, images=%s, isDebug=%t", scanPath, resolutionFilePath, images, isDebug)

	//1. extract files
	filesWithImages, _ := files.ExtractFiles(scanPath)

	//2. extract images from files
	imagesFromFiles := files.ExtractImagesFromFiles(filesWithImages)

	//3. merge all images
	imagesToAnalyze := files.MergeImages(toImageModels(images), imagesFromFiles)

	//4. get images resolution
	resolutionResult, err := syftExtractor.AnalyzeImages(imagesToAnalyze)
	if err != nil {
		log.Fatal("Could not analyze images", err)
	}

	//5. save to resolution file path
	err = files.SaveObjectToFile(resolutionFilePath, resolutionResult)
	if err != nil {
		log.Fatal("Could not save resolution result", err)
	}
}

func toImageModels(images []string) []files.ImageModel {
	var imageNames []files.ImageModel

	for _, image := range images {
		imageNames = append(imageNames, files.ImageModel{
			Name:   image,
			Origin: files.UserInput,
			Path:   files.NoFilePath,
		})
	}

	return imageNames
}
