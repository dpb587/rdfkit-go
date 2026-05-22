package jsonldrdfio

import (
	"fmt"
	"strings"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldcontent"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/iri/rdfacontext"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type encoder struct{}

var _ rdfiotypes.EncoderManager = encoder{}

func NewEncoder() rdfiotypes.EncoderManager {
	return encoder{}
}

func (encoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return jsonldcontent.TypeIdentifier
}

func (e encoder) NewEncoderParams() rdfiotypes.Params {
	return &encoderParams{}
}

func (e encoder) NewEncoder(ww rdfiotypes.Writer, opts rdfiotypes.EncoderOptions) (*rdfiotypes.EncoderHandle, error) {
	params := &encoderParams{}

	err := rdfiotypes.LoadAndApplyParams(params, opts.Params...)
	if err != nil {
		return nil, fmt.Errorf("params: %v", err)
	}

	options := jsonld.EncoderConfig{}

	if bnStringProvider := rdfiotypes.PropagateDecoderPipeBlankNodeStringProvider(opts.DecoderPipe); bnStringProvider != nil {
		options = options.SetBlankNodeStringProvider(bnStringProvider)
	}

	if *params.Buffered {
		options = options.SetBuffered(true)
	}

	if params.Pretty != nil && *params.Pretty {
		options = options.SetIndent("", "\t")
	}

	if *params.IrisUseBase && len(opts.BaseIRI) > 0 {
		options = options.SetBase(string(opts.BaseIRI))
	}

	{
		var prefixes iri.PrefixMappingList

		for _, prefix := range params.IrisUsePrefixes {
			if prefix == "rdfa-context" {
				prefixes = rdfacontext.AppendWidelyUsedInitialContext(prefixes)

				continue
			} else if prefix == "none" {
				prefixes = nil

				continue
			}

			prefixSplit := strings.SplitN(prefix, ":", 2)
			if len(prefixSplit) != 2 {
				return nil, fmt.Errorf("flag[prefixes]: invalid prefix format")
			}

			prefixes = append(prefixes, iri.PrefixMapping{
				Prefix:   prefixSplit[0],
				Expanded: prefixSplit[1],
			})
		}

		if len(prefixes) > 0 {
			options = options.SetPrefixes(prefixes)
		}
	}

	allOptions, err := rdfiotypes.PatchGenericOptions([]jsonld.EncoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	encoder, err := jsonld.NewEncoder(ww, allOptions...)
	if err != nil {
		return nil, err
	}

	return &rdfiotypes.EncoderHandle{
		Writer:  ww,
		Encoder: encoder,
	}, nil
}
