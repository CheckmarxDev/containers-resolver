package main

import (
	"github.com/CheckmarxDev/containers-resolver/pkg/containerResolver"
	"log"
)

const defaultImage = "rabbitmq:3"
const defaultImage2 = "golang:1.21.5-alpine3.18"

func main() {

	scanPath := "./test_files/empty-folder"
	//scanPath := "./test_files/withHelmInZip.zip"
	//scanPath := "./test_files/withDockerInTar.tar.gz"
	//scanPath := "path-to-local-dir"

	resultPath := "./test_files"
	//resultPath := "./test_files/helm-results"
	//resultPath := "./test_files/tar-results"
	//resultPath := "./test_files/dir-results"

	err := containersResolver.Resolve(scanPath, resultPath, []string{}, true)
	if err != nil {
		log.Println("Could not resolve containers", err)
	}
}
