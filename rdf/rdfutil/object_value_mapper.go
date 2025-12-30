package rdfutil

import "github.com/dpb587/rdfkit-go/rdf"

type ObjectValueMapperFunc func(lexicalForm string) (rdf.ObjectValue, error)

func CoalesceObjectValue(lexicalForm string, mappers ...ObjectValueMapperFunc) (rdf.ObjectValue, error) {
	for _, mapper := range mappers {
		v, err := mapper(lexicalForm)
		if err == nil {
			return v, nil
		}
	}

	return nil, rdf.ErrLiteralLexicalFormNotValid
}
