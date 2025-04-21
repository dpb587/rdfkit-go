package encoding

import (
	"errors"
)

var (
	// TODO encodingutil
	ExceedsMaxUnicodePointErr = errors.New("exceeds maximum unicode code point")
)
