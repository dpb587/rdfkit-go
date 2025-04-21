package rdfxml

import (
	"encoding/xml"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectxml-go/inspectxml"
)

type unifiedAttr struct {
	xml.Attr

	Metadata *inspectxml.TokenAttributeMetadata
}

// TODO callers should have better nil checks
var emptyAttrMetadata = &inspectxml.TokenAttributeMetadata{}

func (d *Decoder) getUnifiedAttributes(t xml.StartElement, tokenMetadata *inspectxml.TokenMetadata) []unifiedAttr {
	attrs := make([]unifiedAttr, len(t.Attr))

	if tokenMetadata == nil {
		for i, attr := range t.Attr {
			attrs[i] = unifiedAttr{
				Attr:     attr,
				Metadata: emptyAttrMetadata,
			}
		}

		return attrs
	}

	for i, attr := range t.Attr {
		attrs[i] = unifiedAttr{
			Attr:     attr,
			Metadata: tokenMetadata.TagAttr[i],
		}

		if attrs[i].Metadata == nil {
			attrs[i].Metadata = emptyAttrMetadata
		}
	}

	return attrs
}

func (d *Decoder) newTokenError(err error, t xml.Token) error {
	if d.tokenMetadata == nil {
		return err
	}

	tokenMetadata, ok := d.tokenMetadata()
	if !ok {
		return err
	}

	return cursorio.OffsetRangeError{
		Err:         err,
		OffsetRange: tokenMetadata.Token,
	}
}

func (d *Decoder) newTokenNameError(err error, t xml.Token) error {
	if d.tokenMetadata == nil {
		return err
	}

	tokenMetadata, ok := d.tokenMetadata()
	if !ok || tokenMetadata.TagName == nil {
		return err
	}

	return cursorio.OffsetRangeError{
		Err:         err,
		OffsetRange: *tokenMetadata.TagName,
	}
}

func (d *Decoder) newTokenAttrError(err error, t unifiedAttr) error {
	if d.tokenMetadata == nil || t.Metadata == nil {
		return err
	}

	return cursorio.OffsetRangeError{
		Err:         err,
		OffsetRange: t.Metadata.Name,
	}
}
