package htmlmicrodata

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type DecoderConfig struct {
	vocabularyResolver VocabularyResolver

	laxContentAttributeUse  *bool
	laxContentAttributeHook func(err DecoderError_LaxContentAttribute)
}

var _ DecoderOption = DecoderConfig{}

func (b DecoderConfig) SetVocabularyResolver(v VocabularyResolver) DecoderConfig {
	b.vocabularyResolver = v

	return b
}

func (b DecoderConfig) SetLaxContentAttribute(use bool, hook func(err DecoderError_LaxContentAttribute)) DecoderConfig {
	b.laxContentAttributeUse = &use
	b.laxContentAttributeHook = hook

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.vocabularyResolver != nil {
		s.vocabularyResolver = b.vocabularyResolver
	}

	if b.laxContentAttributeUse != nil && b.laxContentAttributeHook != nil {
		s.laxContentAttributeUse = b.laxContentAttributeUse
		s.laxContentAttributeHook = b.laxContentAttributeHook
	}
}

func (b DecoderConfig) newDecoder(doc *encodinghtml.Document) (*Decoder, error) {
	docProfile := doc.GetInfo()

	w := &Decoder{
		doc:                doc,
		captureOffsets:     docProfile.HasNodeMetadata,
		vocabularyResolver: LiteralVocabularyResolver,
		buildTextOffsets:   encodingutil.BuildTextOffsetsNil,
		statementIdx:       -1,
	}

	if len(docProfile.BaseURL) > 0 {
		docBaseURL, err := iriutil.ParseIRI(docProfile.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("parse document url: %v", err)
		}

		w.docBaseURL = docBaseURL
	}

	if docProfile.HasNodeMetadata {
		w.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	if b.vocabularyResolver != nil {
		w.vocabularyResolver = b.vocabularyResolver
	}

	if b.laxContentAttributeUse != nil && b.laxContentAttributeHook != nil {
		w.laxContentAttributeUse = *b.laxContentAttributeUse
		w.laxContentAttributeHook = b.laxContentAttributeHook
		w.laxContentAttribute = *b.laxContentAttributeUse || b.laxContentAttributeHook != nil
	}

	return w, nil
}
