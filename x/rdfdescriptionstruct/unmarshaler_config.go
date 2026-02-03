package rdfdescriptionstruct

import "github.com/dpb587/rdfkit-go/iri"

// UnmarshalerOption is an option for configuring an Unmarshaler.
type UnmarshalerOption interface {
	apply(*Unmarshaler)
}

// UnmarshalerConfig provides configuration for an Unmarshaler.
type UnmarshalerConfig struct {
	prefixes iri.PrefixMappingList
}

// SetPrefixes overrides the default RDFa Initial Context (Widely Used) prefixes.
func (c UnmarshalerConfig) SetPrefixes(prefixes iri.PrefixMappingList) UnmarshalerConfig {
	c.prefixes = prefixes
	return c
}

// apply implements UnmarshalerOption.
func (c UnmarshalerConfig) apply(u *Unmarshaler) {
	u.prefixes = iri.NewPrefixManager(c.prefixes)
}

// NewUnmarshalerConfig creates a new UnmarshalerConfig with defaults.
func NewUnmarshalerConfig() UnmarshalerConfig {
	return UnmarshalerConfig{}
}
