package xsdtype

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/objecttypes"
)

var durationValidRE = regexp.MustCompile(`^(-)?P(((\d*(\.\d*)?)Y)?((\d*(\.\d*)?)M)?((\d*(\.\d*)?)D)?)?(T((\d*(\.\d*)?)H)?((\d*(\.\d*)?)M)?((\d*(\.\d*)?)S)?)?$`)

type Duration struct {
	Years    float64
	Months   float64
	Days     float64
	Hours    float64
	Minutes  float64
	Seconds  float64
	Negative bool
}

var _ objecttypes.Value = Duration{}

func MapDuration(lexicalForm string) (Duration, error) {
	vMatch := durationValidRE.FindStringSubmatch(xsdutil.WhiteSpaceCollapse(lexicalForm))
	if vMatch == nil {
		return Duration{}, rdf.ErrLiteralLexicalFormNotValid
	}

	var err error
	var l Duration

	if vMatch[1] == "-" {
		l.Negative = true
	}

	if len(vMatch[4]) > 0 {
		l.Years, err = strconv.ParseFloat(vMatch[4], 64)
		if err != nil {
			return Duration{}, fmt.Errorf("%w: year: %v", rdf.ErrLiteralLexicalFormNotValid, err)
		}
	}

	if len(vMatch[7]) > 0 {
		l.Months, err = strconv.ParseFloat(vMatch[7], 64)
		if err != nil {
			return Duration{}, fmt.Errorf("%w: month: %v", rdf.ErrLiteralLexicalFormNotValid, err)
		}
	}

	if len(vMatch[10]) > 0 {
		l.Days, err = strconv.ParseFloat(vMatch[10], 64)
		if err != nil {
			return Duration{}, fmt.Errorf("%w: day: %v", rdf.ErrLiteralLexicalFormNotValid, err)
		}
	}

	if len(vMatch[14]) > 0 {
		l.Hours, err = strconv.ParseFloat(vMatch[14], 64)
		if err != nil {
			return Duration{}, fmt.Errorf("%w: hour: %v", rdf.ErrLiteralLexicalFormNotValid, err)
		}
	}

	if len(vMatch[17]) > 0 {
		l.Minutes, err = strconv.ParseFloat(vMatch[17], 64)
		if err != nil {
			return Duration{}, fmt.Errorf("%w: minute: %v", rdf.ErrLiteralLexicalFormNotValid, err)
		}
	}

	if len(vMatch[20]) > 0 {
		l.Seconds, err = strconv.ParseFloat(vMatch[20], 64)
		if err != nil {
			return Duration{}, fmt.Errorf("%w: second: %v", rdf.ErrLiteralLexicalFormNotValid, err)
		}
	}

	return l, nil
}

func (v Duration) AsObjectValue() rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    xsdiri.Duration_Datatype,
		LexicalForm: v.AsLexicalForm(),
	}
}

func (Duration) TermKind() rdf.TermKind {
	return rdf.TermKindLiteral
}

func (v Duration) TermEquals(t rdf.Term) bool {
	tLiteral, ok := t.(rdf.Literal)
	if !ok {
		return false
	} else if tLiteral.Datatype != xsdiri.Duration_Datatype {
		return false
	}

	return v.AsLexicalForm() == tLiteral.LexicalForm
}

func (l Duration) AsLexicalForm() string {
	out := &strings.Builder{}

	if l.Negative {
		out.WriteByte('-')
	}

	out.WriteByte('P')

	if l.Years > 0 {
		out.WriteString(strconv.FormatFloat(l.Years, 'f', -1, 64))
		out.WriteByte('Y')
	}

	if l.Months > 0 {
		out.WriteString(strconv.FormatFloat(l.Months, 'f', -1, 64))
		out.WriteByte('M')
	}

	if l.Days > 0 {
		out.WriteString(strconv.FormatFloat(l.Days, 'f', -1, 64))
		out.WriteByte('D')
	}

	if l.Hours > 0 || l.Minutes > 0 || l.Seconds > 0 {
		out.WriteByte('T')

		if l.Hours > 0 {
			out.WriteString(strconv.FormatFloat(l.Hours, 'f', -1, 64))
			out.WriteByte('H')
		}

		if l.Minutes > 0 {
			out.WriteString(strconv.FormatFloat(l.Minutes, 'f', -1, 64))
			out.WriteByte('M')
		}

		if l.Seconds > 0 {
			out.WriteString(strconv.FormatFloat(l.Seconds, 'f', -1, 64))
			out.WriteByte('S')
		}
	}

	return out.String()
}
