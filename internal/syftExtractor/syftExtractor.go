package syftExtractor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/linux"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	"log"
	"regexp"
	"strings"
)

func AnalyzeImages(images []types.ImageModel) (*ContainerResolution, error) {

	containerResolution := &ContainerResolution{
		ContainerImages:   []ContainerImage{},
		ContainerPackages: []ContainerPackage{},
	}

	for _, imageModel := range images {
		tmpResolution, err := analyzeImage(imageModel)

		if err != nil {
			log.Printf("Could not analyze image: %s", imageModel)
			continue
		}

		containerResolution.ContainerImages = append(containerResolution.ContainerImages, tmpResolution.ContainerImages...)
		containerResolution.ContainerPackages = append(containerResolution.ContainerPackages, tmpResolution.ContainerPackages...)
		log.Printf("Successfully analyzed image: %s", imageModel)
	}
	return containerResolution, nil
}

func analyzeImage(imageModel types.ImageModel) (*ContainerResolution, error) {

	log.Printf("image is %s, origin: %s, file path: %s", imageModel.Name, imageModel.Origin, imageModel.Path)

	imageSource, s, err := analyzeImageUsingSyft(imageModel.Name)
	if err != nil {
		return nil, err
	}

	result := transformSBOMToContainerResolution(*s, imageSource, imageModel.Name, imageModel.Path, imageModel.Origin)

	return &result, nil
}

func analyzeImageUsingSyft(imageId string) (*source.StereoscopeImageSource, *sbom.SBOM, error) {
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
		log.Printf("Could not pull image: %s. err: %s", imageId, err.Error())
		return nil, nil, fmt.Errorf("could not pull image. %s", err.Error())
	}

	s := getSBOM(imageSource)

	return imageSource, &s, nil
}

func getSBOM(src source.Source) sbom.SBOM {
	s, err := syft.CreateSBOM(context.Background(), src, nil)
	if err != nil {
		panic(err)
	}

	return *s
}

func transformSBOMToContainerResolution(sbom sbom.SBOM, imageSource *source.StereoscopeImageSource, imageId, imagePath, imageOrigin string) ContainerResolution {

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

	extractImage(sbom, imageSource.ID(), imageId, imagePath, imageOrigin, sourceMetadata, imageNameAndTag, &result)
	extractImagePackages(imageId, imageSource.ID(), sbom, &result)

	return result
}

func extractImage(s sbom.SBOM, imageHash artifact.ID, imageId, imagePath, imageOrigin string, sourceMetadata source.StereoscopeImageSourceMetadata, imageNameAndTag []string, result *ContainerResolution) {

	history := extractHistory(sourceMetadata)
	layerIds := extractLayerIds(history)

	images := ContainerImage{
		ImageName:    imageNameAndTag[0],
		ImageTag:     imageNameAndTag[1],
		ImagePath:    imagePath,
		Distribution: getDistro(s.Artifacts.LinuxDistribution),
		ImageHash:    string(imageHash),
		ImageId:      imageId,
		ImageOrigin:  imageOrigin,
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
			Distribution:  getDistro(s.Artifacts.LinuxDistribution),
			Type:          containerPackage.Type.PackageURLType(),
			SourceName:    "",
			SourceVersion: "",
			Licenses:      extractPackageLicenses(containerPackage),
			LayerIds:      extractPackageLayerIds(containerPackage.Locations),
		})
	}

	result.ContainerPackages = packages
}

func getDistro(release *linux.Release) string {
	if release == nil {
		return types.NoFilePath
	}
	return release.PrettyName
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
