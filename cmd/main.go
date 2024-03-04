package main

import (
	"github.com/Checkmarx-Containers/containers-resolver/pkg/containerResolver"
	"os"
)

const defaultImage = "alpine:3.14.0"
const defaultImage2 = "nginx:latest"

func main() {

	i := imageReference()
	i2 := defaultImage2
	//containersResolver.Resolve("/Users/danielgreenspan/Desktop/containers-worker.zip", "containers-resolution-compressed-zip.json", []string{i, i2}, false)
	//containersResolver.Resolve("/Users/danielgreenspan/Desktop/worker.tar.gz", "containers-resolution-compressed-tar.json", []string{i, i2}, false)

	currentPath, _ := os.Getwd() // Get the current working directory (use it as the scan path and resolution folder path)
	containersResolver.Resolve(currentPath, currentPath, []string{i, i2}, false)
}

func imageReference() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return defaultImage
}
