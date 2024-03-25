# containers-resolve
This Go module simplifies the process of analyzing images by providing tools to extract images from various file formats and resolve the software packages within them. It enables users to gain insights into the contents of Docker images, facilitating tasks such as vulnerability assessments and software inventory management. With support for debugging and flexible extraction methods, it's a valuable resource for developers, DevOps engineers, and security professionals working with containerized environments.


## Supported File Types for Package Analysis

This module supports scanning and analyzing the following types of files to extract Docker images and resolve their associated packages:

- **Dockerfile**: Dockerfiles are text documents that contain all the commands a user could call on the command line to assemble an image. This module can parse Dockerfiles to identify image dependencies and extract Docker images specified within them.

- **Docker Compose Files**: Docker Compose is a tool used to define and run multi-container Docker applications. This module can process Docker Compose YAML files to extract Docker images referenced within them, enabling analysis of the entire application stack.

- **Helm Charts**: Helm is a package manager for Kubernetes that provides a way to define, install, and manage Kubernetes applications. Helm charts, which are YAML files, define the structure and configuration of Kubernetes resources. This module can parse Helm charts to extract Docker images used in deploying Kubernetes applications.

By supporting these file types, this module offers versatility in scanning and analyzing various sources of Docker images, catering to different deployment scenarios and containerization strategies.

## Installation

You can install this module using Go modules. 

```bash
go get github.com/your-username/your-repo
```


# Usage

To use this module, import it in your Go code and call the `Resolve` function with appropriate parameters:

```go
import (
    "fmt"
    "github.com/your-username/your-repo/module"
)

func main() {
    err := module.Resolve(scanPath, resolutionFolderPath, images, isDebug)
    if err != nil {
        fmt.Printf("Error occurred: %v\n", err)
        // Handle error
        return
    }
    fmt.Println("Resolution successful!")
}
```

# Functionality Overview
Resolve
The Resolve function scans Docker images and extracts their SBOMs.

``` bash
func Resolve(scanPath string, resolutionFolderPath string, images []string, isDebug bool) error
```

## Parameters
- scanPath: Path to the directory containing source code (can be a folder path or a path to zipped files).
- resolutionFolderPath: Path to the folder where resolution results will be saved.
- images: List of image names to scan besids what is found in the `scanPath` folder.
- isDebug: Boolean flag indicating whether logging debug mode is enabled.

## Returns
#### error: Returns an error if any operation fails.

## License
This module is distributed under the MIT License.

## Contributing
Contributions are welcome! Please follow the guidelines outlined in the CONTRIBUTING.md file.

## Additional Information

- Requires Go version 1.21.8 or later.
- Compatible with Linux, macOS, and Windows.
