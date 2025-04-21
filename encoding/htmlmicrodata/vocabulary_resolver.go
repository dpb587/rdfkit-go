package htmlmicrodata

import (
	"net/url"
)

type VocabularyResolver interface {
	ResolveMicrodataProperty(itemtypes []string, itemprop string) (string, error)
}

//

type literalVocabularyResolver struct{}

var LiteralVocabularyResolver VocabularyResolver = &literalVocabularyResolver{}

func (r *literalVocabularyResolver) ResolveMicrodataProperty(itemtypes []string, itemprop string) (string, error) {
	return itemprop, nil
}

//

type itemtypeVocabularyResolver struct{}

var ItemtypeVocabularyResolver VocabularyResolver = &itemtypeVocabularyResolver{}

func (r *itemtypeVocabularyResolver) ResolveMicrodataProperty(itemtypes []string, itemprop string) (string, error) {
	if len(itemtypes) == 0 {
		return itemprop, nil
	}

	tt, err := url.Parse(itemtypes[0])
	if err != nil {
		return itemprop, nil
	}

	ttv, err := tt.Parse(itemprop)
	if err != nil {
		return itemprop, nil
	}

	return ttv.String(), nil
}
