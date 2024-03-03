package files

type FileImages struct {
	Dockerfile    []FilePath
	DockerCompose []FilePath
}

type FilePath struct {
	FullPath     string
	RelativePath string
}

type ImageModel struct {
	Name   string
	Origin string
	Path   string
}

const (
	UserInput               = "UserInput"
	DockerFileOrigin        = "Dockerfile"
	DockerComposeFileOrigin = "DockerCompose"
	NoFilePath              = "NONE"
)
