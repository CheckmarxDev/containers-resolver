package main

import (
	"github.com/Checkmarx-Containers/containers-resolver/pkg/containerResolver"
	"log"
)

const defaultImage = "debian:12"
const defaultImage2 = "nginx:latest"

func main() {

	//scanPath := "./test_files/withDockerInZip.zip"
	scanPath := "./test_files/withHelmInZip.zip"
	//scanPath := "./test_files/withDockerInTar.tar.gz"
	//scanPath := "path-to-local-dir"

	//resultPath := "./test_files/zip-results"
	resultPath := "./test_files/helm-results"
	//resultPath := "./test_files/tar-results"
	//resultPath := "./test_files/dir-results"

	err := containersResolver.Resolve(scanPath, resultPath, []string{defaultImage, defaultImage2}, true)
	if err != nil {
		log.Println("Could not resolve containers", err)
	}
}
