package nquads

import (
	"io"

	"github.com/dpb587/rdfkit-go/encoding/nquads/internal"
	"github.com/dpb587/rdfkit-go/rdf"
)

func WriteIRI(w io.Writer, t rdf.IRI, ascii bool) (int, error) {
	var uchar4, uchar8 int

	tr := []rune(t)

	for idx := 0; idx < len(tr); idx++ {
		switch iriMustEscapeRune(tr[idx], ascii) {
		case iriRuneEscapeUCHAR4:
			uchar4++
		case iriRuneEscapeUCHAR8:
			uchar8++
		}
	}

	if uchar4 == 0 && uchar8 == 0 {
		return w.Write([]byte("<" + string(t) + ">"))
	}

	buf := make([]rune, len(tr)+uchar4*5+uchar8*9)
	widx := 0

	for ridx := 0; ridx < len(tr); ridx++ {
		rr := tr[ridx]

		switch iriMustEscapeRune(rr, ascii) {
		case iriRuneEscapeUCHAR4:
			buf[widx] = '\\'
			buf[widx+1] = 'u'
			buf[widx+2] = rune(internal.HexUpper[rr&0xf000>>12])
			buf[widx+3] = rune(internal.HexUpper[rr&0x0f00>>8])
			buf[widx+4] = rune(internal.HexUpper[rr&0x00f0>>4])
			buf[widx+5] = rune(internal.HexUpper[rr&0x000f])
			widx += 6
		case iriRuneEscapeUCHAR8:
			buf[widx] = '\\'
			buf[widx+1] = 'U'
			buf[widx+2] = rune(internal.HexUpper[rr&0x70000000>>28])
			buf[widx+3] = rune(internal.HexUpper[rr&0x0f000000>>24])
			buf[widx+4] = rune(internal.HexUpper[rr&0x00f00000>>20])
			buf[widx+5] = rune(internal.HexUpper[rr&0x000f0000>>16])
			buf[widx+6] = rune(internal.HexUpper[rr&0x0000f000>>12])
			buf[widx+7] = rune(internal.HexUpper[rr&0x00000f00>>8])
			buf[widx+8] = rune(internal.HexUpper[rr&0x000000f0>>4])
			buf[widx+9] = rune(internal.HexUpper[rr&0x0000000f])
			widx += 10
		default:
			buf[widx] = rr
			widx++
		}
	}

	return w.Write([]byte("<" + string(buf) + ">"))
}

type iriRuneEscapeMode uint

const (
	iriRuneEscapeNone iriRuneEscapeMode = iota
	iriRuneEscapeUCHAR4
	iriRuneEscapeUCHAR8
)

func iriMustEscapeRune(r rune, ascii bool) iriRuneEscapeMode {
	if r >= 0x00 && r <= 0x20 {
		return iriRuneEscapeUCHAR4
	} else if r == '<' || r == '>' || r == '"' || r == '{' || r == '}' || r == '|' || r == '^' || r == '`' || r == '\\' {
		return iriRuneEscapeUCHAR4
	} else if ascii {
		if r > 0xffff {
			return iriRuneEscapeUCHAR8
		} else if r > 0xff {
			return iriRuneEscapeUCHAR4
		}
	}

	return iriRuneEscapeNone
}
