package turtle

import (
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal"
)

func format_PN_LOCAL(v string) string {
	var usePercent, escapeEsc int

	tr := []rune(v)

	for idx := 0; idx < len(tr); idx++ {
		switch prefixLocalNameMustEscapeRune(tr[idx]) {
		case prefixLocalNameRuneEscapePERCENT:
			usePercent++
		case prefixLocalNameRuneEscapeESC:
			escapeEsc++
		}
	}

	if usePercent == 0 && escapeEsc == 0 {
		return v
	}

	buf := make([]rune, len(tr)+usePercent*2+escapeEsc)
	widx := 0

	for ridx := 0; ridx < len(tr); ridx++ {
		rr := tr[ridx]

		switch prefixLocalNameMustEscapeRune(rr) {
		case prefixLocalNameRuneEscapePERCENT:
			buf[widx] = '%'
			buf[widx+1] = rune(cursorioutil.HexUpper[rr&0x00f0>>4])
			buf[widx+2] = rune(cursorioutil.HexUpper[rr&0x000f])
			widx += 3
		case prefixLocalNameRuneEscapeESC:
			buf[widx] = '\\'
			buf[widx+1] = rr
			widx += 2
		default:
			buf[widx] = rr
			widx++
		}
	}

	return string(buf)
}

type prefixLocalNameRuneEscapeMode uint

const (
	prefixLocalNameRuneEscapeNone prefixLocalNameRuneEscapeMode = iota
	prefixLocalNameRuneEscapePERCENT
	prefixLocalNameRuneEscapeESC
)

// technically there are different cases for first/last/other runes of the sequence
// should fix eventually with more complicated, index-based conditionals
func prefixLocalNameMustEscapeRune(r rune) prefixLocalNameRuneEscapeMode {
	switch r {
	case '_', '~', '.', '-', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', '/', '?', '#', '@', '%':
		return prefixLocalNameRuneEscapeESC
	case ':':
		return prefixLocalNameRuneEscapeNone
	}

	if internal.IsRune_PN_CHARS(r) {
		return prefixLocalNameRuneEscapeNone
	}

	return prefixLocalNameRuneEscapePERCENT
}
