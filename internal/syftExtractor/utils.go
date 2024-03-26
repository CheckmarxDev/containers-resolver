package syftExtractor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/format"
	"github.com/anchore/syft/syft/format/syftjson"
	"github.com/anchore/syft/syft/linux"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	"regexp"
	"strings"
)

func analyzeImage(l *logger.Logger, imageModel types.ImageModel) (*ContainerResolution, error) {

	l.Debug("image is %s, found in file paths: %s", imageModel.Name, GetImageLocationsPathsString(imageModel))

	imageSource, s, err := analyzeImageUsingSyft(l, imageModel.Name)
	if err != nil {
		return nil, err
	}

	result := transformSBOMToContainerResolution(l, *s, imageSource, imageModel)

	return &result, nil
}

func analyzeImageUsingSyft(l *logger.Logger, imageId string) (*source.StereoscopeImageSource, *sbom.SBOM, error) {
	platform, err := image.NewPlatform("linux/amd64")
	if err != nil {
		l.Error("could not create platform object", err)
		return nil, nil, err
	}

	imageSource, err := source.NewFromStereoscopeImage(
		source.StereoscopeImageConfig{
			Reference: imageId,
			From:      image.DockerDaemonSource,
			Platform:  platform,
		},
	)
	if err != nil {
		l.Error("Could not pull image: %s. err: %+v", imageId, err)
		return nil, nil, err
	}

	s, err := getSBOM(imageSource, true)
	if err != nil {
		l.Error("Could get image SBOM. image: %s. err: %+v", imageId, err)
		return nil, nil, err
	}
	return imageSource, &s, nil
}

func getSBOM(src source.Source, saveToFile bool) (sbom.SBOM, error) {
	s, err := syft.CreateSBOM(context.Background(), src, nil)
	if err != nil {
		return sbom.SBOM{}, err
	}

	if saveToFile {
		formatSBOM(*s)
	}
	return *s, nil
}

func formatSBOM(s sbom.SBOM) []byte {
	bytes, err := format.Encode(s, syftjson.NewFormatEncoder())
	if err != nil {
		panic(err)
	}
	return bytes
}

func transformSBOMToContainerResolution(l *logger.Logger, s sbom.SBOM, imageSource *source.StereoscopeImageSource, imageModel types.ImageModel) ContainerResolution {

	imageNameAndTag := strings.Split(imageModel.Name, ":")

	imageResult := ContainerResolution{
		ContainerImage:    ContainerImage{},
		ContainerPackages: []ContainerPackage{},
	}
	var sourceMetadata source.StereoscopeImageSourceMetadata
	var ok bool

	if sourceMetadata, ok = s.Source.Metadata.(source.StereoscopeImageSourceMetadata); !ok {
		l.Warn("Value is not StereoscopeImageSourceMetadata - can not analyze")
		return imageResult
	}

	distro := getDistro(s.Artifacts.LinuxDistribution)

	extractImage(distro, imageSource.ID(), imageModel, sourceMetadata, imageNameAndTag, &imageResult)
	extractImagePackages(s.Artifacts.Packages, distro, &imageResult)

	return imageResult
}

func extractImage(distro string, imageHash artifact.ID, imageModel types.ImageModel, sourceMetadata source.StereoscopeImageSourceMetadata, imageNameAndTag []string, result *ContainerResolution) {

	history := extractHistory(sourceMetadata)
	layerIds := extractLayerIds(history)

	result.ContainerImage = ContainerImage{
		ImageName:      imageNameAndTag[0],
		ImageTag:       imageNameAndTag[1],
		Distribution:   distro,
		ImageHash:      string(imageHash),
		ImageId:        imageModel.Name,
		Layers:         layerIds,
		History:        history,
		ImageLocations: getImageLocations(imageModel.ImageLocations),
	}
}

func extractImagePackages(packages *pkg.Collection, distro string, result *ContainerResolution) {

	var containerPackages []ContainerPackage

	for containerPackage := range packages.Enumerate() {

		sourceName, sourceVersion := getPackageRelationships(containerPackage)

		containerPackages = append(containerPackages, ContainerPackage{
			Name:          containerPackage.Name,
			Version:       containerPackage.Version,
			Distribution:  distro,
			Type:          packageTypeToPackageManager(containerPackage.Type),
			SourceName:    sourceName,
			SourceVersion: sourceVersion,
			Licenses:      extractPackageLicenses(containerPackage),
			LayerIds:      extractPackageLayerIds(containerPackage.Locations),
		})
	}

	result.ContainerPackages = containerPackages
}

func getPackageRelationships(containerPackage pkg.Package) (string, string) {

	if apkMeta, ok := containerPackage.Metadata.(pkg.ApkDBEntry); ok {
		return getApkSource(containerPackage, apkMeta)
	}
	if debMeta, ok := containerPackage.Metadata.(pkg.DpkgDBEntry); ok {
		return getDebSource(containerPackage, debMeta)
	}
	if rpmMeta, ok := containerPackage.Metadata.(pkg.RpmDBEntry); ok {
		return getRpmSource(containerPackage, rpmMeta)
	}
	return "", ""
}

func getApkSource(containerPackage pkg.Package, apkMeta pkg.ApkDBEntry) (string, string) {
	if apkMeta.OriginPackage == "" || apkMeta.OriginPackage == containerPackage.Name {
		return "", ""
	}
	if apkMeta.Version == "" {
		return apkMeta.OriginPackage, containerPackage.Version
	}
	return apkMeta.OriginPackage, apkMeta.Version
}

func getDebSource(containerPackage pkg.Package, debMeta pkg.DpkgDBEntry) (string, string) {
	if debMeta.Source == "" || debMeta.Source == containerPackage.Name {
		return "", ""
	}
	if debMeta.SourceVersion == "" {
		return debMeta.Source, containerPackage.Version
	}
	return debMeta.Source, debMeta.SourceVersion
}

func getRpmSource(containerPackage pkg.Package, rpmMeta pkg.RpmDBEntry) (string, string) {
	if rpmMeta.SourceRpm == "" || rpmMeta.SourceRpm == containerPackage.Name {
		return "", ""
	}
	if rpmMeta.SourceRpm == "" {
		return rpmMeta.SourceRpm, containerPackage.Version
	}
	return rpmMeta.SourceRpm, rpmMeta.Version
}

func getDistro(release *linux.Release) string {
	if release == nil || release.ID == "" || release.VersionID == "" {
		return types.NoFilePath
	}
	return fmt.Sprintf("%s:%s", release.ID, trimPatchVersion(release.VersionID))
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

func trimPatchVersion(versionID string) string {
	re := regexp.MustCompile(`^(\d+\.\d+)`)
	matches := re.FindStringSubmatch(versionID)

	if len(matches) > 1 {
		return matches[1]
	}
	return versionID
}

func getImageLocations(imageLocations []types.ImageLocation) []ImageLocation {
	var slice []ImageLocation
	for _, location := range imageLocations {
		slice = append(slice, ImageLocation{
			Origin: location.Origin,
			Path:   location.Path,
		})
	}
	return slice
}

func GetImageLocationsPathsString(imgModel types.ImageModel) string {
	var paths []string
	for _, location := range imgModel.ImageLocations {
		paths = append(paths, location.Path)
	}
	return strings.Join(paths, ", ")
}
