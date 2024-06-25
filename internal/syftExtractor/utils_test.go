package syftExtractor

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/anchore/syft/syft/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractPackageName(t *testing.T) {
	tests := []struct {
		name     string
		pack     pkg.Package
		expected string
	}{
		{
			name: "ExtractName for supported package type",
			pack: pkg.Package{
				Name: "examplePackage",
				Type: "maven",
				Metadata: pkg.JavaArchive{
					PomProperties: &pkg.JavaPomProperties{
						GroupID: "exampleGroup",
					},
				},
			},
			expected: "exampleGroup:examplePackage",
		},
		{
			name: "No extraction for unsupported package type",
			pack: pkg.Package{
				Name: "examplePackage",
				Type: "unsupportedType",
			},
			expected: "examplePackage",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := extractPackageName(test.pack)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetGroupId(t *testing.T) {
	tests := []struct {
		name     string
		metadata interface{}
		expected string
	}{
		{
			name:     "JavaMetadata with GroupID",
			metadata: pkg.JavaArchive{PomProperties: &pkg.JavaPomProperties{GroupID: "testGroup"}},
			expected: "testGroup",
		},
		{
			name:     "JavaMetadata without GroupID",
			metadata: pkg.JavaArchive{},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getGroupId(test.metadata)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestOutputFormat(t *testing.T) {
	expected := "testGroup:testPackage"
	result := outputFormat("testGroup", "testPackage")
	assert.Equal(t, expected, result)
}

func TestExtractName(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		purl        string
		groupId     string
		expected    string
	}{
		{
			name:        "ExtractName with GroupID and PURL",
			packageName: "testPackage",
			purl:        "https://example.com/testGroup/testPackage",
			groupId:     "testGroup",
			expected:    "testGroup:testPackage",
		},
		{
			name:        "ExtractName with PURL but no GroupID",
			packageName: "testPackage",
			purl:        "https://example.com/invalid",
			groupId:     "",
			expected:    "testPackage",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := extractName(test.packageName, test.purl, test.groupId)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetSyftArtifactsWithoutUnsupportedTypesDuplications(t *testing.T) {
	mockLogger := logger.NewLogger(false)

	pkg1 := pkg.Package{Name: "test1", Version: "1.0", Type: pkg.JavaPkg}
	pkg2 := pkg.Package{Name: "test2", Version: "2.0", Type: pkg.ApkPkg}
	pkg3 := pkg.Package{Name: "test1", Version: "1.0", Type: pkg.UnknownPkg}
	pkg4 := pkg.Package{Name: "test4", Version: "2.0", Type: pkg.JavaPkg}

	collection := pkg.NewCollection(pkg1, pkg2, pkg3, pkg4)

	expected := []pkg.Package{pkg1, pkg2, pkg4}

	result := getSyftArtifactsWithoutUnsupportedTypesDuplications(mockLogger, collection)

	assert.Equal(t, len(expected), len(result))
}
