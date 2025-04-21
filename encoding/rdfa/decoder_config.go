package rdfa

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type DecoderConfig struct {
	htmlProcessingProfile *HtmlProcessingProfile
	defaultVocabulary     *string
	blankNodeStringMapper blanknodeutil.StringMapper
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

func (b DecoderConfig) SetBlankNodeStringMapper(v blanknodeutil.StringMapper) DecoderConfig {
	b.blankNodeStringMapper = v

	return b
}

func (b DecoderConfig) apply(s *DecoderConfig) {
	if b.htmlProcessingProfile != nil {
		s.htmlProcessingProfile = b.htmlProcessingProfile
	}

	if b.defaultVocabulary != nil {
		s.defaultVocabulary = b.defaultVocabulary
	}

	if b.blankNodeStringMapper != nil {
		s.blankNodeStringMapper = b.blankNodeStringMapper
	}
}

func (b DecoderConfig) newDecoder(doc *encodinghtml.Document) (*Decoder, error) {
	docProfile := doc.GetInfo()

	w := &Decoder{
		doc:                   doc,
		captureOffsets:        docProfile.HasNodeMetadata,
		defaultVocabulary:     b.defaultVocabulary,
		blankNodeStringMapper: b.blankNodeStringMapper,
		buildTextOffsets:      encodingutil.BuildTextOffsetsNil,
		statementIdx:          -1,
	}

	if len(docProfile.BaseURL) > 0 {
		docBaseURL, err := iriutil.ParseIRI(docProfile.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("parse document url: %v", err)
		}

		w.docBaseURL = docBaseURL
	}

	if b.htmlProcessingProfile != nil {
		w.htmlProcessingProfile = *b.htmlProcessingProfile
	}

	if w.blankNodeStringMapper == nil {
		w.blankNodeStringMapper = blanknodeutil.NewStringMapper()
	}

	if docProfile.HasNodeMetadata {
		w.buildTextOffsets = encodingutil.BuildTextOffsetsValue
	}

	return w, nil
}
