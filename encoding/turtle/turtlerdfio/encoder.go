package turtlerdfio

import (
	"context"
	"fmt"
	"strings"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/encoding/turtle/turtlecontent"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil/rdfacontext"
	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionutil"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
)

type encoder struct{}

var _ rdfiotypes.EncoderManager = encoder{}

func NewEncoder() rdfiotypes.EncoderManager {
	return encoder{}
}

func (encoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return turtlecontent.TypeIdentifier
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

	options := turtle.EncoderConfig{}

	if *params.Buffered {
		options = options.SetBuffered(true)
	}

	if *params.IrisUseBase == true {
		options = options.SetBase(string(opts.BaseIRI))
	}

	{
		var prefixes iriutil.PrefixMappingList

		for _, prefix := range params.IrisUsePrefixes {
			if prefix == "rdfa-context" {
				prefixes = append(prefixes, rdfacontext.WidelyUsedInitialContext()...)

				continue
			} else if prefix == "none" {
				prefixes = nil

				continue
			}

			prefixSplit := strings.SplitN(prefix, ":", 2)
			if len(prefixSplit) != 2 {
				return nil, fmt.Errorf("flag[prefixes]: invalid prefix format")
			}

			prefixes = append(prefixes, iriutil.PrefixMapping{
				Prefix:   prefixSplit[0],
				Expanded: rdf.IRI(prefixSplit[1]),
			})
		}

		if len(prefixes) > 0 {
			options = options.SetPrefixes(iriutil.NewPrefixMap(prefixes...))
		}
	}

	allOptions, err := rdfiotypes.PatchGenericOptions([]turtle.EncoderOption{options}, opts.Patcher)
	if err != nil {
		return nil, err
	}

	encoder, err := turtle.NewEncoder(ww, allOptions...)
	if err != nil {
		return nil, err
	}

	var wrappedEncoder encoding.Encoder = encoder

	if *params.Resources {
		wrappedEncoder = rdfdescriptionutil.NewBufferedTriplesEncoder(
			context.Background(),
			encoder,
			true,
		)
	}

	return &rdfiotypes.EncoderHandle{
		Writer:  ww,
		Encoder: wrappedEncoder,
	}, nil
}
