package rdfxml

import "github.com/dpb587/rdfkit-go/encoding"

type DocumentResource struct{}

var _ encoding.ContainerResource = &DocumentResource{}

func (*DocumentResource) ContainerResourceString() string {
	return "document"
}
