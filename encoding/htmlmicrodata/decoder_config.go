package htmlmicrodata

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type DecoderConfig struct {
	vocabularyResolver VocabularyResolver
}

var _ DecoderOption = DecoderConfig{}

func (b DecoderConfig) SetVocabularyResolver(v VocabularyResolver) DecoderConfig {
	b.vocabularyResolver = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.vocabularyResolver != nil {
		s.vocabularyResolver = b.vocabularyResolver
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

	return w, nil
}
