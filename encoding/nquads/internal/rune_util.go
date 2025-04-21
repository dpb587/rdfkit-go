package internal

// PN_CHARS ::= PN_CHARS_U | '-' | [0-9] | #x00B7 | [#x0300-#x036F] | [#x203F-#x2040]
func IsRune_PN_CHARS(r rune) bool {
	switch {
	case false,
		IsRune_PN_CHARS_U(r),
		r == '-',
		'0' <= r && r <= '9',
		r == 0x00B7,
		0x0300 <= r && r <= 0x036F,
		0x203F <= r && r <= 0x2040:
		return true
	}

	return false
}

func IsRune_PN_CHARS_U(r rune) bool {
	switch {
	case false,
		IsRune_PN_CHARS_BASE(r),
		r == '_',
		r == ':':
		return true
	}

	return false
}

// [A-Z] | [a-z] | [#x00C0-#x00D6] | [#x00D8-#x00F6] | [#x00F8-#x02FF] | [#x0370-#x037D] | [#x037F-#x1FFF] | [#x200C-#x200D] | [#x2070-#x218F] | [#x2C00-#x2FEF] | [#x3001-#xD7FF] | [#xF900-#xFDCF] | [#xFDF0-#xFFFD] | [#x10000-#xEFFFF]
func IsRune_PN_CHARS_BASE(r rune) bool {
	switch {
	case false,
		'A' <= r && r <= 'Z',
		'a' <= r && r <= 'z',
		0x00C0 <= r && r <= 0x00D6,
		0x00D8 <= r && r <= 0x00F6,
		0x00F8 <= r && r <= 0x02FF,
		0x0370 <= r && r <= 0x037D,
		0x037F <= r && r <= 0x1FFF,
		0x200C <= r && r <= 0x200D,
		0x2070 <= r && r <= 0x218F,
		0x2C00 <= r && r <= 0x2FEF,
		0x3001 <= r && r <= 0xD7FF,
		0xF900 <= r && r <= 0xFDCF,
		0xFDF0 <= r && r <= 0xFFFD,
		0x10000 <= r && r <= 0xEFFFF:
		return true
	}

	return false
}
