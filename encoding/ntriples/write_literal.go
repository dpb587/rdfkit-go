package ntriples

import (
	"io"

	"github.com/dpb587/rdfkit-go/encoding/ntriples/internal"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

func WriteLiteral(w io.Writer, t rdf.Literal, ascii bool) (int, error) {
	var wlen int

	var echar, uchar4, uchar8 int

	tr := []rune(t.LexicalForm)

	for idx := 0; idx < len(tr); idx++ {
		switch literalStringMustEscapeRune(tr[idx], ascii) {
		case stringLiteralQuoteRuneEscapeECHAR:
			echar++
		case stringLiteralQuoteRuneEscapeUCHAR4:
			uchar4++
		case stringLiteralQuoteRuneEscapeUCHAR8:
			uchar8++
		}
	}

	if echar == 0 && uchar4 == 0 && uchar8 == 0 {
		wwlen, err := w.Write([]byte(`"` + string(t.LexicalForm) + `"`))
		if err != nil {
			return wwlen, err
		}

		wlen += wwlen
	} else {
		buf := make([]rune, len(tr)+echar+uchar4*5+uchar8*9)
		widx := 0

		for ridx := 0; ridx < len(tr); ridx++ {
			rr := tr[ridx]

			switch literalStringMustEscapeRune(rr, ascii) {
			case stringLiteralQuoteRuneEscapeECHAR:
				buf[widx] = '\\'

				switch rr {
				case '\t':
					buf[widx+1] = 't'
				case '\b':
					buf[widx+1] = 'b'
				case '\n':
					buf[widx+1] = 'n'
				case '\r':
					buf[widx+1] = 'r'
				case '\f':
					buf[widx+1] = 'f'
				case '"':
					buf[widx+1] = '"'
				case '\'':
					buf[widx+1] = '\''
				case '\\':
					buf[widx+1] = '\\'
				}

				widx += 2
			case stringLiteralQuoteRuneEscapeUCHAR4:
				buf[widx] = '\\'
				buf[widx+1] = 'u'
				buf[widx+2] = rune(internal.HexUpper[rr&0xf000>>12])
				buf[widx+3] = rune(internal.HexUpper[rr&0x0f00>>8])
				buf[widx+4] = rune(internal.HexUpper[rr&0x00f0>>4])
				buf[widx+5] = rune(internal.HexUpper[rr&0x000f])
				widx += 6
			case stringLiteralQuoteRuneEscapeUCHAR8:
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

		wwlen, err := w.Write([]byte(`"` + string(buf) + `"`))
		if err != nil {
			return wwlen, err
		}

		wlen += wwlen
	}

	if t.Datatype == xsdiri.String_Datatype {
		return wlen, nil
	} else if t.Datatype == rdfiri.LangString_Datatype {
		wwlen, err := w.Write([]byte("@" + t.Tags[rdf.LanguageLiteralTag]))

		return wlen + wwlen, err
	}

	wwlen, err := w.Write([]byte("^^"))
	if err != nil {
		return wlen + wwlen, err
	}

	wwlen, err = WriteIRI(w, t.Datatype, ascii)

	return wlen + wwlen, err
}

type stringLiteralQuoteRuneEscapeMode uint

const (
	stringLiteralQuoteRuneEscapeNone stringLiteralQuoteRuneEscapeMode = iota
	// stringLiteralQuoteRuneEscapeInvalid
	stringLiteralQuoteRuneEscapeECHAR
	stringLiteralQuoteRuneEscapeUCHAR4
	stringLiteralQuoteRuneEscapeUCHAR8
)

func literalStringMustEscapeRune(r rune, ascii bool) stringLiteralQuoteRuneEscapeMode {
	if r == 0x0022 || r == 0x005C || r == 0x000A || r == 0x000D {
		return stringLiteralQuoteRuneEscapeECHAR
	} else if ascii {
		if r > 0xffff {
			return stringLiteralQuoteRuneEscapeUCHAR8
		} else if r > 0xff {
			return stringLiteralQuoteRuneEscapeUCHAR4
		}
	}

	return stringLiteralQuoteRuneEscapeNone
}
