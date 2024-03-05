package extractors

import (
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/types"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"log"
	"regexp"

	"path/filepath"
	"strings"
)

func ExtractImagesFromHelmFiles(helmCharts []types.HelmChartInfo) ([]types.ImageModel, error) {

	var imagesFromHelmDirectories []types.ImageModel
	for _, h := range helmCharts {
		renderedTemplates, err := generateRenderedTemplates(h)
		if err != nil {
			log.Println("Could not get images from helm dir", h.Directory)
		}

		images, err := extractImageInfo(renderedTemplates)
		if err != nil {
			return nil, err
		}

		log.Printf("Found images in helm directory: %s, images: %v", h.Directory, images)
		imagesFromHelmDirectories = append(imagesFromHelmDirectories, images...)
	}

	return imagesFromHelmDirectories, nil
}

func generateRenderedTemplates(c types.HelmChartInfo) (string, error) {
	actionConfig := new(action.Configuration)

	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ReleaseName = "temp-release"
	client.ClientOnly = true

	chartPath, err := filepath.Abs(c.Directory)
	if err != nil {
		return "", err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return "", err
	}

	release, err := client.Run(chart, nil)
	if err != nil {
		return "", err
	}

	return release.Manifest, nil
}

func extractImageInfo(yamlString string) ([]types.ImageModel, error) {
	sections := strings.Split(yamlString, "---")

	var imageInfoList []types.ImageModel

	for _, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}

		var microservice types.Microservice
		err := yaml.Unmarshal([]byte(section), &microservice)
		if err != nil {
			return nil, err
		}

		s, _ := extractSource(section)
		n, _ := extractImageName(microservice)

		v := types.ImageModel{
			Name:   n,
			Origin: types.HelmFileOrigin,
			Path:   s,
		}

		imageInfoList = append(imageInfoList, v)
	}

	return imageInfoList, nil
}

func extractImageName(microservice types.Microservice) (string, error) {
	var imageName string
	if microservice.Spec.Image.Registry != "" {
		imageName += microservice.Spec.Image.Registry + "/"
	}
	imageName += microservice.Spec.Image.Name + ":"
	imageName += microservice.Spec.Image.Tag

	return imageName, nil
}

func extractSource(yamlBlock string) (string, error) {
	sourceRegex := regexp.MustCompile(`#\s*Source:\s*([^\n]+)`)
	match := sourceRegex.FindStringSubmatch(yamlBlock)

	if len(match) != 2 {
		return "", fmt.Errorf("source not found in YAML block")
	}

	source := strings.TrimSpace(match[1])
	return source, nil
}
