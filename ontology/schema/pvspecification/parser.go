package pvspecification

import (
	"regexp"
)

func ParseText(v string) ShorthandText {
	var lex = v
	var res ShorthandText

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

		res.Attributes = append(res.Attributes, spec)
	}

	return res
}
