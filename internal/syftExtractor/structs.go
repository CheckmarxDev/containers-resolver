package syftExtractor

type ImageConfig struct {
	History []HistoryConfig `json:"history"`
	Rootfs  RootfsConfig    `json:"rootfs"`
}

type HistoryConfig struct {
	CreatedBy  string `json:"created_by"`
	EmptyLayer bool   `json:"empty_layer"`
}

type RootfsConfig struct {
	DiffIds []string `json:"diff_ids"`
}

type ContainerResolution struct {
	ContainerImages   []ContainerImage
	ContainerPackages []ContainerPackage
}

type ContainerImage struct {
	ImageName    string
	ImageTag     string
	ImagePath    string
	Distribution string
	ImageHash    string
	ImageId      string
	ImageOrigin  string
	Layers       []string
	History      []Layer
}

type ContainerPackage struct {
	ImageId       string
	ImageHash     string
	Name          string
	Version       string
	Distribution  string
	Type          string
	SourceName    string
	SourceVersion string
	Licenses      []string
	LayerIds      []string
}

type Layer struct {
	Order   int
	Size    int64
	LayerId string
	Command string
}
