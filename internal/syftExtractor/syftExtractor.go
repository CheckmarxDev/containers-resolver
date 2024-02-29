package syftExtractor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	"log"
	"regexp"
	"strings"
)

func AnalyzeImages(images []string) (*ContainerResolution, error) {

	containerResolution := &ContainerResolution{
		ContainerImages:   []ContainerImage{},
		ContainerPackages: []ContainerPackage{},
	}

	for _, imageName := range images {
		tmpResolution, err := analyzeImage(imageName)

		if err != nil {
			log.Printf("Could not analyze image: %s", imageName)
			continue
		}

		containerResolution.ContainerImages = append(containerResolution.ContainerImages, tmpResolution.ContainerImages...)
		containerResolution.ContainerPackages = append(containerResolution.ContainerPackages, tmpResolution.ContainerPackages...)

	}

	return containerResolution, nil
}

func analyzeImage(imageId string) (*ContainerResolution, error) {

	log.Printf("image is %s", imageId)

	imageSource, s := analyzeImageUsingSyft(imageId)

	result := transformSBOMToContainerResolution(s, imageSource, imageId, "some-path/Dockerfile")

	return &result, nil
}

func analyzeImageUsingSyft(imageId string) (*source.StereoscopeImageSource, sbom.SBOM) {
	platform, err := image.NewPlatform("linux/amd64")
	if err != nil {
		panic(err)
	}

	imageSource, err := source.NewFromStereoscopeImage(
		source.StereoscopeImageConfig{
			Reference: imageId,
			From:      image.DockerDaemonSource,
			Platform:  platform,
		},
	)
	if err != nil {
		panic(err)
	}

	s := getSBOM(imageSource)

	return imageSource, s
}

func getSBOM(src source.Source) sbom.SBOM {
	s, err := syft.CreateSBOM(context.Background(), src, nil)
	if err != nil {
		panic(err)
	}

	return *s
}

func transformSBOMToContainerResolution(sbom sbom.SBOM, imageSource *source.StereoscopeImageSource, imageId, imagePath string) ContainerResolution {

	imageNameAndTag := strings.Split(imageId, ":")

	result := ContainerResolution{
		ContainerImages:   []ContainerImage{},
		ContainerPackages: []ContainerPackage{},
	}
	var sourceMetadata source.StereoscopeImageSourceMetadata
	var ok bool

	if sourceMetadata, ok = sbom.Source.Metadata.(source.StereoscopeImageSourceMetadata); !ok {
		fmt.Println("Value is not StereoscopeImageSourceMetadata - can not analyze")
		return result
	}

	extractImage(sbom, imageSource.ID(), imageId, imagePath, sourceMetadata, imageNameAndTag, &result)
	extractImagePackages(imageId, imageSource.ID(), sbom, &result)

	return result
}

func extractImage(s sbom.SBOM, imageHash artifact.ID, imageId string, imagePath string, sourceMetadata source.StereoscopeImageSourceMetadata, imageNameAndTag []string, result *ContainerResolution) {

	history := extractHistory(sourceMetadata)
	layerIds := extractLayerIds(history)

	images := ContainerImage{
		ImageName:    imageNameAndTag[0],
		ImageTag:     imageNameAndTag[1],
		ImagePath:    imagePath,
		Distribution: s.Artifacts.LinuxDistribution.PrettyName,
		ImageHash:    string(imageHash),
		ImageId:      imageId,
		ImageOrigin:  "user-input",
		Layers:       layerIds,
		History:      history,
	}

	result.ContainerImages = append(result.ContainerImages, images)
}

func extractImagePackages(imageId string, imageHash artifact.ID, s sbom.SBOM, result *ContainerResolution) {

	var packages []ContainerPackage

	for containerPackage := range s.Artifacts.Packages.Enumerate() {
		packages = append(packages, ContainerPackage{
			ImageId:       imageId,
			ImageHash:     string(imageHash),
			Name:          containerPackage.Name,
			Version:       containerPackage.Version,
			Distribution:  s.Artifacts.LinuxDistribution.PrettyName,
			Type:          containerPackage.Type.PackageURLType(),
			SourceName:    "",
			SourceVersion: "",
			Licenses:      extractPackageLicenses(containerPackage),
			LayerIds:      extractPackageLayerIds(containerPackage.Locations),
		})
	}

	result.ContainerPackages = packages
}

func extractPackageLayerIds(locations file.LocationSet) []string {
	var layerIds []string
	for _, l := range locations.ToSlice() {
		layerIds = append(layerIds, removeSha256(l.FileSystemID))
	}
	return layerIds
}

func extractPackageLicenses(p pkg.Package) []string {
	var licenses []string
	for _, l := range p.Licenses.ToSlice() {
		licenses = append(licenses, l.Value)
	}
	return licenses
}

func extractLayerIds(layers []Layer) []string {
	var layerIds []string

	for _, layer := range layers {
		if layer.LayerId != "" {
			layerIds = append(layerIds, layer.LayerId)
		}
	}

	return layerIds
}

func extractHistory(sourceMetadata source.StereoscopeImageSourceMetadata) []Layer {
	imageConfig := decodeBase64ToJson(sourceMetadata.RawConfig)
	j := 0

	var res []Layer
	for i := 0; i < len(imageConfig.History); i++ {
		isLayerEmpty := imageConfig.History[i].EmptyLayer
		var layerID string
		if !isLayerEmpty {
			layerID = removeSha256(imageConfig.Rootfs.DiffIds[j])
		}

		res = append(res, Layer{
			Order:   i,
			Size:    getSize(layerID, sourceMetadata.Layers),
			LayerId: layerID,
			Command: imageConfig.History[i].CreatedBy,
		})

		if !isLayerEmpty {
			j++
		}
	}
	return res
}

func decodeBase64ToJson(base64Bytes []byte) ImageConfig {
	var imageConfig ImageConfig
	err := json.Unmarshal(base64Bytes, &imageConfig)
	if err != nil {
		return ImageConfig{}
	}
	return imageConfig
}

func removeSha256(str string) string {
	if strings.TrimSpace(str) == "" {
		return str
	}
	return regexp.MustCompile(`^sha256:`).ReplaceAllString(str, "")
}

func getSize(layerId string, layers []source.StereoscopeLayerMetadata) int64 {
	for _, layer := range layers {
		if removeSha256(layer.Digest) == layerId {
			return layer.Size
		}
	}
	return 0
}
