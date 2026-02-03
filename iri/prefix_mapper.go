package iri

type PrefixMapper interface {
	CompactPrefix(v string) (PrefixReference, bool)
	ExpandPrefix(pr PrefixReference) (string, bool)
}
