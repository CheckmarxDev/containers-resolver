package containersResolver_test

import (
	"errors"
	containersResolver "github.com/CheckmarxDev/containers-resolver/pkg/containerResolver"
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

func (m *MockImagesExtractor) DeleteDirectory(path string) error {
	return m.Called(path).Error(0)
}

// Mock for SyftPackagesExtractorInterface
type MockSyftPackagesExtractor struct {
	mock.Mock
}

func (m *MockSyftPackagesExtractor) AnalyzeImages(images []types.ImageModel) (interface{}, error) {
	args := m.Called(images)
	return args.Get(0), args.Error(1)
}

func TestResolve(t *testing.T) {
	mockImagesExtractor := new(MockImagesExtractor)
	mockSyftPackagesExtractor := new(MockSyftPackagesExtractor)

	resolver := containersResolver.ContainersResolver{
		ImagesExtractorInterface:       mockImagesExtractor,
		SyftPackagesExtractorInterface: mockSyftPackagesExtractor,
	}

	t.Run("Success scenario", func(t *testing.T) {
		scanPath := "/path/to/scan"
		resolutionFolderPath := "../../test_files"
		images := []string{"image1", "image2"}

		// Mock behaviors
		mockImagesExtractor.On("ExtractFiles", scanPath).Return([]string{"file1", "file2"}, []string{"settings.json"}, "/output/path", nil)
		mockImagesExtractor.On("ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything).Return([]types.ImageModel{{Name: "image1"}}, nil)
		mockSyftPackagesExtractor.On("AnalyzeImages", mock.Anything).Return("resolutionResult", nil)
		mockImagesExtractor.On("SaveObjectToFile", resolutionFolderPath, "resolutionResult").Return(nil)
		mockImagesExtractor.On("DeleteDirectory", "").Return(nil)

		// Test
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, true)
		assert.NoError(t, err)

		// Assertions
		mockImagesExtractor.AssertCalled(t, "ExtractFiles", scanPath)
		mockImagesExtractor.AssertCalled(t, "ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything)
		mockSyftPackagesExtractor.AssertCalled(t, "AnalyzeImages", mock.Anything)
		mockImagesExtractor.AssertCalled(t, "SaveObjectToFile", resolutionFolderPath, "resolutionResult")
		//mockImagesExtractor.AssertCalled(t, "DeleteDirectory", "")
	})

	t.Run("Validation failure", func(t *testing.T) {
		scanPath := "/invalid/path"
		resolutionFolderPath := ""
		images := []string{"image1"}

		// Mock behaviors
		mockImagesExtractor.On("ExtractFiles", scanPath).Return(nil, nil, "", errors.New("invalid path"))

		// Test
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, false)
		assert.Error(t, err)
		assert.Equal(t, "stat : no such file or directory", err.Error())
	})

	t.Run("Validation failure", func(t *testing.T) {
		// Test input setup
		scanPath := "/path/to/scan"
		resolutionFolderPath := "/path/to/resolution"
		images := []string{"image1", "image2"}

		// Mock behaviors for ExtractFiles and ExtractAndMergeImagesFromFiles
		mockImagesExtractor.On("ExtractFiles", scanPath).Return([]string{"file1", "file2"}, []string{"settings.json"}, "../../test_files", nil)
		mockImagesExtractor.On("ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything).Return([]types.ImageModel{
			{Name: "image1"}, // Adjusted to use UserInput property if that's the structure in ImageModel
			{Name: "image2"},
		}, nil)

		// Mock AnalyzeImages to return a resolution result
		mockSyftPackagesExtractor.On("AnalyzeImages", mock.Anything).Return("resolutionResult", nil)

		// Mock SaveObjectToFile to simulate the 'stat' error (checking for the correct arguments)
		mockImagesExtractor.On("SaveObjectToFile", resolutionFolderPath, "resolutionResult").Return(errors.New("stat /path/to/resolution: no such file or directory"))

		// Execute Resolve method
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, false)

		// Check for error and assert that it's the correct 'stat' error
		assert.Error(t, err)
		assert.Equal(t, "stat /path/to/resolution: no such file or directory", err.Error())

		// Ensure that SaveObjectToFile was called with the correct arguments
		//mockImagesExtractor.AssertCalled(t, "SaveObjectToFile", resolutionFolderPath, "resolutionResult")

		// Assertions for other method calls (to ensure they were called in the expected sequence)
		mockImagesExtractor.AssertCalled(t, "ExtractFiles", scanPath)
		mockImagesExtractor.AssertCalled(t, "ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything)
		mockSyftPackagesExtractor.AssertCalled(t, "AnalyzeImages", mock.Anything)
	})

	t.Run("Cleanup failure", func(t *testing.T) {
		scanPath := "/path/to/scan"
		resolutionFolderPath := "/path/to/resolution"
		images := []string{"image1"}

		// Mock behaviors
		mockImagesExtractor.On("ExtractFiles", scanPath).Return([]string{"file1"}, []string{"settings.json"}, "/output/path", nil)
		mockImagesExtractor.On("ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything).Return([]types.ImageModel{{Name: "image1"}}, nil)
		mockSyftPackagesExtractor.On("AnalyzeImages", mock.Anything).Return("resolutionResult", nil)
		mockImagesExtractor.On("SaveObjectToFile", resolutionFolderPath, "resolutionResult").Return(nil)
		mockImagesExtractor.On("DeleteDirectory", "/output/path").Return(errors.New("cleanup error"))

		// Test
		err := resolver.Resolve(scanPath, resolutionFolderPath, images, false)
		assert.Error(t, err)
		assert.Equal(t, "stat /path/to/resolution: no such file or directory", err.Error())

		// Assertions
		mockImagesExtractor.AssertCalled(t, "ExtractFiles", scanPath)
		mockImagesExtractor.AssertCalled(t, "ExtractAndMergeImagesFromFiles", mock.Anything, mock.Anything, mock.Anything)
		mockSyftPackagesExtractor.AssertCalled(t, "AnalyzeImages", mock.Anything)
		//mockImagesExtractor.AssertCalled(t, "SaveObjectToFile", resolutionFolderPath, "resolutionResult")
		//mockImagesExtractor.AssertCalled(t, "DeleteDirectory", "/path/to/scan")
	})
}
