package ntriples

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
)

func (t *Decoder) newOffsetError(err error, readUncomitted, readIgnored []rune) error {
	werr := cursorio.OffsetError{
		Err: err,
	}

	if t.doc == nil {
		o := t.buf.GetByteOffset()

		if len(readIgnored) > 0 {
			o -= cursorio.ByteOffset(cursorioutil.RunesBytes(readIgnored))
		}

		werr.Offset = o
	} else {
		if len(readUncomitted) > 0 {
			werr.Offset = *t.uncommittedTextOffset(readUncomitted)
		} else {
			werr.Offset = t.doc.GetTextOffset()
		}
	}

	return werr
}

func (t *Decoder) getTextOffset() *cursorio.TextOffset {
	if t.doc == nil {
		return nil
	}

	v := t.doc.GetTextOffset()

	return &v
}

func (t *Decoder) commit(runes []rune) {
	if t.doc == nil {
		return
	}

	t.doc.WriteRunes(runes)
}

func (t *Decoder) commitForTextOffsetRange(runes []rune) *cursorio.TextOffsetRange {
	if t.doc == nil {
		return nil
	}

	v := t.doc.WriteRunesForOffsetRange(runes)

	return &v
}

func (t *Decoder) uncommittedTextOffset(runes []rune) *cursorio.TextOffset {
	if t.doc == nil {
		return nil
	}

	clone := t.doc.Clone()
	v := clone.WriteRunesForOffset(runes)

	return &v
}

func (t *Decoder) uncommittedTextOffsetRange(prefixIgnored, runes []rune) *cursorio.TextOffsetRange {
	if t.doc == nil {
		return nil
	}

	clone := t.doc.Clone()

	if len(prefixIgnored) > 0 {
		clone.WriteRunes(prefixIgnored)
	}

	v := clone.WriteRunesForOffsetRange(runes)

	return &v
}
