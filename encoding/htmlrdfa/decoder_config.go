package htmlrdfa

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/iri/rdfacontext"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type DecoderConfig struct {
	htmlProcessingProfile *HtmlProcessingProfile
	defaultVocabulary     *string
	defaultPrefixes       iri.PrefixMappingList
	bnStringFactory       blanknodes.StringFactory
}

var _ DecoderOption = DecoderConfig{}

func (b DecoderConfig) SetHtmlProcessingProfile(v HtmlProcessingProfile) DecoderConfig {
	b.htmlProcessingProfile = &v

	return b
}

func (b DecoderConfig) SetDefaultVocabulary(v string) DecoderConfig {
	b.defaultVocabulary = &v

	return b
}

func (b DecoderConfig) SetDefaultPrefixes(v iri.PrefixMappingList) DecoderConfig {
	b.defaultPrefixes = v

	return b
}

func (b DecoderConfig) SetBlankNodeStringFactory(v blanknodes.StringFactory) DecoderConfig {
	b.bnStringFactory = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.htmlProcessingProfile != nil {
		s.htmlProcessingProfile = b.htmlProcessingProfile
	}

	if b.defaultVocabulary != nil {
		s.defaultVocabulary = b.defaultVocabulary
	}

	if b.defaultPrefixes != nil {
		s.defaultPrefixes = b.defaultPrefixes
	}

	if b.bnStringFactory != nil {
		s.bnStringFactory = b.bnStringFactory
	}
}

var emptyURL = (func() *iri.ParsedIRI {
	iri, err := iri.ParseIRI("")
	if err != nil {
		panic(fmt.Sprintf("failed to parse empty IRI: %v", err))
	}

	return iri
})()

func (b DecoderConfig) newDecoder(doc *encodinghtml.Document) (*Decoder, error) {
	docProfile := doc.GetInfo()

	w := &Decoder{
		doc:               doc,
		captureOffsets:    docProfile.HasNodeMetadata,
		defaultVocabulary: b.defaultVocabulary,
		bnStringFactory:   b.bnStringFactory,
		buildTextOffsets:  encodingutil.BuildTextOffsetsNil,
		statementsIdx:     -1,
	}

	if len(b.defaultPrefixes) > 0 {
		w.defaultPrefixes = iri.NewPrefixManager(b.defaultPrefixes)
	} else {
		w.defaultPrefixes = rdfacontext.NewWidelyUsedInitialContext()
	}

	if len(docProfile.BaseURL) > 0 {
		docBaseURL, err := iri.ParseIRI(docProfile.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("parse document url: %v", err)
		}

		// test[0117.html] fragment must be dropped; not documented in specs?
		docBaseURL.DropFragment()

		w.docBaseURL = docBaseURL
	} else {
		w.docBaseURL = emptyURL
	}

	if b.htmlProcessingProfile != nil {
		w.htmlProcessingProfile = *b.htmlProcessingProfile
	}

	if w.bnStringFactory == nil {
		w.bnStringFactory = blanknodes.NewStringFactory()
	}

	if docProfile.HasNodeMetadata {
		w.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	return w, nil
}
