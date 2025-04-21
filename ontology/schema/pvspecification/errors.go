package pvspecification

import (
	"fmt"
)

type UnknownShorthandTextAttributeNameError struct {
	Name string
}

func (e UnknownShorthandTextAttributeNameError) Error() string {
	return fmt.Sprintf("unknown shorthand attribute name: %s", e.Name)
}
