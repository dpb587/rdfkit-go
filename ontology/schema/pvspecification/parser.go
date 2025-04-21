package pvspecification

import (
	"regexp"
)

type ParseTextResult struct {
	Data ShorthandText

	// UnknownAttributeNameIndices contains any indices of Data.Attributes which did not have a valid mapping for name.
	UnknownAttributeNameIndices []int
}

// Err will be non-nil if there were any unknown attribute names in the parsed text.
func (r ParseTextResult) Err() error {
	if len(r.UnknownAttributeNameIndices) > 0 {
		return UnknownShorthandTextAttributeNameError{
			Name: r.Data.Attributes[r.UnknownAttributeNameIndices[0]].Name,
		}
	}

	return nil
}

func ParseText(v string) ParseTextResult {
	var res ParseTextResult
	var lex = v

	for len(lex) > 0 {
		nextMatch := regexp.MustCompile(`\s*([^\s=]+)(=([^\s]+)?)?`).FindStringSubmatchIndex(lex)
		if nextMatch == nil {
			break
		}

		spec := ShorthandTextAttribute{
			Name: lex[nextMatch[2]:nextMatch[3]],
		}

		if nextMatch[4] > 0 {
			var shorthandValue string

			if nextMatch[6] > 0 {
				shorthandValue = lex[nextMatch[6]:nextMatch[7]]
			}

			spec.Value = &shorthandValue
		}

		lex = lex[nextMatch[1]:]

		res.Data.Attributes = append(res.Data.Attributes, spec)
	}

	return res
}
