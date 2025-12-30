package rdfdescriptionstruct

import (
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

// defaultUnmarshaler is the internal unmarshaler used by the package-level functions.
var defaultUnmarshaler = NewUnmarshaler()

// UnmarshalResource unmarshals an rdfdescription.Resource into a struct.
// The struct fields must have rdf tags indicating how to map RDF data.
func UnmarshalResource(from rdfdescription.Resource, to any) error {
	return defaultUnmarshaler.UnmarshalResource(from, to)
}

// Unmarshal unmarshals an rdfdescription.Resource into a struct using a ResourceListBuilder.
// The builder is required for recursive unmarshaling of resource-type object values.
func Unmarshal(builder *rdfdescription.ResourceListBuilder, from rdfdescription.Resource, to any) error {
	return defaultUnmarshaler.Unmarshal(builder, from, to)
}
