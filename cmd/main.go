package main

import (
	containersResolver "github.com/CheckmarxDev/containers-resolver/pkg/containerResolver"
	"log"
)

const defaultImage = "manuelbcd/vulnapp:latest"
const defaultImage2 = "library/debian:10"

func main() {

	//scanPath := "./test_files/empty-folder"
	//scanPath := "./test_files/withHelmInZip.zip"
	//scanPath := "./test_files/withDockerInTar.tar.gz"
	scanPath := "/Users/danielgreenspan/GolandProjects/containers-resolver/test_files/imageExtraction/dockerfiles/test-envs"

	resultPath := "./test_files"
	//resultPath := "./test_files/helm-results"
	//resultPath := "./test_files/tar-results"
	//resultPath := "./test_files/dir-results"

	err := containersResolver.Resolve(scanPath, resultPath, []string{}, true)
	if err != nil {
		log.Println("Could not resolve containers", err)
	}
}
