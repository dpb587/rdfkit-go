package encoding

// TODO wip

// containers are an experimental concept and likely to change
type ContainerProvider interface {
	GetEncodingContainer() (Container, bool)
}

type Container struct {
	Resource ContainerResource
}

type ContainerResource interface {
	ContainerResourceString() string
}
