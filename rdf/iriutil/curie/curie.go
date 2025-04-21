package curie

import (
	"fmt"
)

type CURIE struct {
	Safe          bool
	DefaultPrefix bool

	Prefix    string
	Reference string
}

func (c CURIE) SafeString() string {
	if c.DefaultPrefix {
		return fmt.Sprintf("[%s]", c.Reference)
	}

	return fmt.Sprintf("[%s:%s]", c.Prefix, c.Reference)
}

func (c CURIE) String() string {
	if c.Safe {
		return c.SafeString()
	} else if c.DefaultPrefix {
		return c.Reference
	}

	return fmt.Sprintf("%s:%s", c.Prefix, c.Reference)
}

//

type CURIEs []CURIE
