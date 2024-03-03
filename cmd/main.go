package main

import (
	"github.com/Checkmarx-Containers/containers-resolver/pkg/abstraction"
	"os"
)

const defaultImage = "alpine:3.19"
const defaultImage2 = "nginx:latest"

func main() {

	i := imageReference()
	i2 := defaultImage2

	containersResolver.Resolve("./app", "containers-resolution.json", []string{i, i2}, false)
}

func imageReference() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return defaultImage
}
