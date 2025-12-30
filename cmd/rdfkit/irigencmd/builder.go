package irigencmd

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/rdfs/rdfsiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

type builder struct {
	base        rdf.IRI
	baseTerm    rdf.IRI
	baseMatcher *regexp.Regexp

	statementsBySubject map[rdf.IRI]*builderSubject
}

func (b *builder) DetectBase() (string, bool) {
	var groups = map[string]int{}

	for _, subject := range b.statementsBySubject {
		if s := strings.SplitN(string(subject.IRI), "#", 2); len(s) == 2 {
			groups[s[0]+"#"]++
		} else {
			s := strings.Split(string(subject.IRI), "/")

			groups[strings.Join(s[:len(s)-1], "/")+"/"]++
		}
	}

	var maxBase string
	var maxCount = 0
	var maxTied bool

	for base, count := range groups {
		if count > maxCount {
			maxBase = base
			maxCount = count
			maxTied = false
		} else if count == maxCount {
			maxTied = true
		}
	}

	if maxTied || maxCount == 0 {
		return "", false
	}

	return maxBase, true
}

func (b *builder) FilterBase(base rdf.IRI) *builder {
	b2 := &builder{
		base:                base,
		statementsBySubject: map[rdf.IRI]*builderSubject{},
	}

	if strings.HasSuffix(string(base), "#") {
		b2.baseMatcher = regexp.MustCompile(fmt.Sprintf("^%s(#|$)", regexp.QuoteMeta(strings.TrimSuffix(string(base), "#"))))
	} else {
		b2.baseMatcher = regexp.MustCompile(fmt.Sprintf("^%s", regexp.QuoteMeta(string(base))))
	}

	for _, subject := range b.statementsBySubject {
		if subject.IRI == base {
			b2.baseTerm = subject.IRI
		}

		if b2.baseMatcher.MatchString(string(subject.IRI)) {
			b2.statementsBySubject[subject.IRI] = subject
		}
	}

	return b2
}

func (b *builder) AddStatement(quad rdf.Quad) {
	sIRI, ok := quad.Triple.Subject.(rdf.IRI)
	if !ok {
		return
	}

	pIRI, ok := quad.Triple.Predicate.(rdf.IRI)
	if !ok {
		return
	}

	if _, ok := b.statementsBySubject[sIRI]; !ok {
		b.statementsBySubject[sIRI] = &builderSubject{
			IRI: sIRI,
		}
	}

	switch pIRI {
	case rdfiri.Type_Property:
		if oIRI, ok := quad.Triple.Object.(rdf.IRI); ok {
			b.statementsBySubject[sIRI].Types = append(b.statementsBySubject[sIRI].Types, oIRI)
		}
	case rdfsiri.Comment_Property:
		if oLiteral, ok := quad.Triple.Object.(rdf.Literal); ok {
			if len(oLiteral.LexicalForm) == 0 {
				return
			}

			switch oLiteral.Datatype {
			case xsdiri.String_Datatype, rdfiri.LangString_Datatype:
				b.statementsBySubject[sIRI].Comments = append(
					b.statementsBySubject[sIRI].Comments,
					oLiteral,
				)
			}
		}
	}
}

func (b *builder) ListDefinedTerms() []rdf.IRI {
	var terms []rdf.IRI

	for _, subject := range b.statementsBySubject {
		if subject.IRI == b.base {
			continue
		} else if len(subject.Types) == 0 {
			continue
		}

		terms = append(terms, subject.IRI)
	}

	slices.SortFunc(terms, func(i, j rdf.IRI) int {
		if b.baseTerm == i {
			return -1
		} else if b.baseTerm == j {
			return 1
		}

		return strings.Compare(b.GetGoIdent(i), b.GetGoIdent(j))
	})

	return terms
}

var replaceFirst = map[byte]string{
	'0': "Zero",
	'1': "One",
	'2': "Two",
	'3': "Three",
	'4': "Four",
	'5': "Five",
	'6': "Six",
	'7': "Seven",
	'8': "Eight",
	'9': "Nine",
}

func (b *builder) GetGoIdent(t rdf.IRI) string {
	baseIdentity := b.baseMatcher.ReplaceAllString(string(t), "")

	typeSuffix := string(b.statementsBySubject[t].Types[0])
	if strings.Contains(string(typeSuffix), "#") {
		typeSuffix = strings.SplitN(string(typeSuffix), "#", 2)[1]
	} else {
		typeSuffixSplit := strings.Split(string(typeSuffix), "/")
		typeSuffix = typeSuffixSplit[len(typeSuffixSplit)-1]
	}

	baseIdentity = b.safeIdent(baseIdentity)
	if len(baseIdentity) == 0 {
		baseIdentity = "Base"
	} else {
		baseIdentity += "_"
	}

	ident := baseIdentity + b.safeIdent(typeSuffix)

	if replacer, ok := replaceFirst[ident[0]]; ok {
		ident = replacer + ident[1:]
	}

	return ident
}

func (b *builder) safeIdent(ident string) string {
	ident = regexp.MustCompile(`(\s+|_)(\w)`).ReplaceAllStringFunc(
		strings.NewReplacer(
			":", "_",
		).Replace(ident),
		func(s string) string {
			return s[0:len(s)-1] + strings.ToUpper(s[len(s)-1:])
		},
	)

	ident = strings.ReplaceAll(ident, " ", "")
	if len(ident) == 0 {
		return ""
	}

	return strings.ToUpper(ident[0:1]) + ident[1:]
}

func (b builder) WriteGoComment(w io.Writer, t rdf.IRI, linePrefix string, skipNL bool) (bool, error) {
	s := b.statementsBySubject[t]
	if s == nil {
		return false, nil
	}

	var comment string

	for _, v := range s.Comments {
		if v.Datatype != rdfiri.LangString_Datatype {
			continue
		} else if tag, ok := v.Tag.(rdf.LanguageLiteralTag); ok {
			if !strings.HasPrefix(tag.Language, "en") {
				continue
			}
		} else {
			continue
		}

		comment = v.LexicalForm

		break
	}

	if len(comment) == 0 {
		for _, v := range s.Comments {
			if v.Datatype != xsdiri.String_Datatype {
				continue
			}

			comment = v.LexicalForm

			break
		}
	}

	if len(comment) == 0 {
		return false, nil
	}

	if !skipNL {
		_, err := fmt.Fprintf(w, "\n")
		if err != nil {
			return true, err
		}
	}

	_, err := fmt.Fprintf(w, "%s", linePrefix)
	if err != nil {
		return true, err
	}

	var lineLength = len(linePrefix)

	for _, field := range strings.Fields(comment) {
		fieldLen := len(bytes.Runes([]byte(field))) + 1
		if lineLength+fieldLen >= maxLineLength {
			_, err = fmt.Fprintf(w, "\n%s", linePrefix)
			if err != nil {
				return true, err
			}

			lineLength = len(linePrefix)
			fieldLen--
		} else if lineLength == len(linePrefix) {
			fieldLen--
		} else {
			_, err = fmt.Fprintf(w, " ")
			if err != nil {
				return true, err
			}
		}

		_, err = fmt.Fprintf(w, "%s", field)
		if err != nil {
			return true, err
		}

		lineLength += fieldLen
	}

	_, err = fmt.Fprintf(w, "\n")

	return true, err
}

type builderSubject struct {
	IRI      rdf.IRI
	Types    []rdf.IRI
	Comments []rdf.Literal
}

var maxLineLength = 120
