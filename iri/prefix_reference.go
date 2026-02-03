package iri

type PrefixReference struct {
	Prefix    string
	Reference string
}

func (pr PrefixReference) String() string {
	return pr.Prefix + ":" + pr.Reference
}
