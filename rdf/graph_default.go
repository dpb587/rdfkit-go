package rdf

// DefaultGraph represents the default graph of a dataset.
const DefaultGraph defaultGraphName = true

type defaultGraphName bool

var _ GraphNameValue = defaultGraphName(false)

func (defaultGraphName) isGraphNameValueBuiltin() {}
