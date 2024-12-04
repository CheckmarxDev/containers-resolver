package containersResolver_test

import (
	"errors"
	containersResolver "github.com/CheckmarxDev/containers-resolver/pkg/containerResolver"
	"github.com/rs/zerolog/log"
	"os"
	"testing"

	"github.com/Checkmarx-Containers/extractor-types/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for ImagesExtractorInterface
type MockImagesExtractor struct {
	mock.Mock
}

func (m *MockImagesExtractor) ExtractFiles(scanPath string) ([]string, []string, string, error) {
	args := m.Called(scanPath)
	return args.Get(0).([]string), args.Get(1).([]string), args.String(2), args.Error(3)
}

func (m *MockImagesExtractor) ExtractAndMergeImagesFromFiles(filesWithImages []string, imageModels []types.ImageModel, settingsFiles []string) ([]types.ImageModel, error) {
	args := m.Called(filesWithImages, imageModels, settingsFiles)
	return args.Get(0).([]types.ImageModel), args.Error(1)
}

func (m *MockImagesExtractor) SaveObjectToFile(filePath string, obj interface{}) error {
	return m.Called(filePath, obj).Error(0)
}

// Mock for SyftPackagesExtractorInterface
type MockSyftPackagesExtractor struct {
	mock.Mock
}

func (m *MockSyftPackagesExtractor) AnalyzeImages(images []types.ImageModel) (interface{}, error) {
	args := m.Called(images)
	return args.Get(0), args.Error(1) // Return first and second values from the mock
}

func createTestFolder(dir string) {
	// Create the directory if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Err(err)
		}
	}
}

func TestResolve(t *testing.T) {
	mockImagesExtractor := new(MockImagesExtractor)
	mockSyftPackagesExtractor := new(MockSyftPackagesExtractor)

	createTestFolder("../../test_files/resolution")

	resolver := containersResolver.ContainersResolver{
		ImagesExtractorInterface:       mockImagesExtractor,
		SyftPackagesExtractorInterface: mockSyftPackagesExtractor,
	}

	t.Run("Success scenario", func(t *testing.T) {
		scanPath := "../../test_files"
		resolutionFolderPath := "../../test_files/resolution"
		images := []string{"image1", "image2"}

		// Mock behaviors
		mockImagesExtractor.On("ExtractFiles", scanPath).Return([]string{"file1", "file2"}, []string{"settings.json"}, "/output/path", nil)
		mockImagesExtractor.On("ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything).Return([]types.ImageModel{{Name: "image1"}}, nil)
		mockSyftPackagesExtractor.On("AnalyzeImages", mock.Anything).Return("resolutionResult", nil)
		mockImagesExtractor.On("SaveObjectToFile", resolutionFolderPath, "resolutionResult").Return(nil)

		// Test
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, true)
		assert.NoError(t, err)

		// Assertions
		mockImagesExtractor.AssertCalled(t, "ExtractFiles", scanPath)
		mockImagesExtractor.AssertCalled(t, "ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything)
		mockSyftPackagesExtractor.AssertCalled(t, "AnalyzeImages", mock.Anything)
		mockImagesExtractor.AssertCalled(t, "SaveObjectToFile", resolutionFolderPath, "resolutionResult")
	})

	t.Run("ScanPath Validation failure", func(t *testing.T) {
		scanPath := "/invalid/path"
		resolutionFolderPath := ""
		images := []string{"image1"}

		// Test
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, false)
		assert.Error(t, err)
		assert.Equal(t, "stat : no such file or directory", err.Error())
	})

	t.Run("ExtractFilesError", func(t *testing.T) {
		mockImagesExtractor.ExpectedCalls = nil
		mockImagesExtractor.Calls = nil

		scanPath := "../../test_files"
		resolutionFolderPath := "../../test_files/resolution"
		images := []string{"image1"}

		// Mock behaviors
		mockImagesExtractor.On("ExtractFiles", scanPath).Return([]string{}, []string{}, "", errors.New("invalid path"))

		// Test
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, false)
		assert.Error(t, err)
		assert.Equal(t, "invalid path", err.Error())
	})

	t.Run("ExtractAndMergeImagesFromFiles_Failure", func(t *testing.T) {
		mockImagesExtractor.ExpectedCalls = nil
		mockImagesExtractor.Calls = nil

		scanPath := "../../test_files"
		resolutionFolderPath := "../../test_files/resolution"
		images := []string{"image1"}

		// Mock behaviors
		mockImagesExtractor.On("ExtractFiles", scanPath).Return([]string{"file1", "file2"}, []string{"settings.json"}, "/output/path", nil)
		mockImagesExtractor.On("ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything).Return([]types.ImageModel{{Name: "image1"}}, errors.New("failed to extract"))

		// Test
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, false)
		assert.Error(t, err)
		assert.Equal(t, "failed to extract", err.Error())
	})

	t.Run("AnalyzeImages_Failure", func(t *testing.T) {
		mockImagesExtractor.ExpectedCalls = nil
		mockImagesExtractor.Calls = nil
		mockSyftPackagesExtractor.ExpectedCalls = nil
		mockSyftPackagesExtractor.Calls = nil

		scanPath := "../../test_files"
		resolutionFolderPath := "../../test_files/resolution"
		images := []string{"image1"}

		// Mock behaviors
		mockImagesExtractor.On("ExtractFiles", scanPath).Return([]string{"file1", "file2"}, []string{"settings.json"}, "/output/path", nil)
		mockImagesExtractor.On("ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything).Return([]types.ImageModel{{Name: "image1"}}, nil)
		mockSyftPackagesExtractor.On("AnalyzeImages", mock.Anything).Return(nil, errors.New("failed to analyze"))
		mockImagesExtractor.On("SaveObjectToFile", resolutionFolderPath, "resolutionResult").Return(nil)

		// Test
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, false)
		assert.Error(t, err)
		assert.Equal(t, "failed to analyze", err.Error())
	})

}
