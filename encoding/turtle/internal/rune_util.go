package internal

// PN_CHARS_BASE ::= [A-Z] | [a-z] | [#xc0-#xd6] | [#xd8-#xf6] | [#xf8-#x2ff] | [#x370-#x37d] | [#x37f-#x1fff] | [#x200c-#x200d] | [#x2070-#x218f] | [#x2c00-#x2fef] | [#x3001-#xd7ff] | [#xf900-#xfdcf] | [#xfdf0-#xfffd] | [#x10000-#xeffff]
func IsRune_PN_CHARS_BASE(r rune) bool {
	switch {
	case 'A' <= r && r <= 'Z':
		return true
	case 'a' <= r && r <= 'z':
		return true
	case 0xC0 <= r && r <= 0xD6:
		return true
	case 0xD8 <= r && r <= 0xF6:
		return true
	case 0xF8 <= r && r <= 0x2FF:
		return true
	case 0x370 <= r && r <= 0x37D:
		return true
	case 0x37F <= r && r <= 0x1FFF:
		return true
	case 0x200C <= r && r <= 0x200D:
		return true
	case 0x2070 <= r && r <= 0x218F:
		return true
	case 0x2C00 <= r && r <= 0x2FEF:
		return true
	case 0x3001 <= r && r <= 0xD7FF:
		return true
	case 0xF900 <= r && r <= 0xFDCF:
		return true
	case 0xFDF0 <= r && r <= 0xFFFD:
		return true
	case 0x10000 <= r && r <= 0xEFFFF:
		return true
	}

	return false
}

// PN_CHARS_U ::= PN_CHARS_BASE | '_'
// PN_CHARS_BASE ::= [A-Z] | [a-z] | [#xc0-#xd6] | [#xd8-#xf6] | [#xf8-#x2ff] | [#x370-#x37d] | [#x37f-#x1fff] | [#x200c-#x200d] | [#x2070-#x218f] | [#x2c00-#x2fef] | [#x3001-#xd7ff] | [#xf900-#xfdcf] | [#xfdf0-#xfffd] | [#x10000-#xeffff]
func IsRune_PN_CHARS_U(r rune) bool {
	switch {
	case 'A' <= r && r <= 'Z':
		return true
	case r == '_':
		return true
	case 'a' <= r && r <= 'z':
		return true
	case 0xC0 <= r && r <= 0xD6:
		return true
	case 0xD8 <= r && r <= 0xF6:
		return true
	case 0xF8 <= r && r <= 0x2FF:
		return true
	case 0x370 <= r && r <= 0x37D:
		return true
	case 0x37F <= r && r <= 0x1FFF:
		return true
	case 0x200C <= r && r <= 0x200D:
		return true
	case 0x2070 <= r && r <= 0x218F:
		return true
	case 0x2C00 <= r && r <= 0x2FEF:
		return true
	case 0x3001 <= r && r <= 0xD7FF:
		return true
	case 0xF900 <= r && r <= 0xFDCF:
		return true
	case 0xFDF0 <= r && r <= 0xFFFD:
		return true
	case 0x10000 <= r && r <= 0xEFFFF:
		return true
	}

	return false
}

// PN_CHARS ::= PN_CHARS_U | '-' | [0-9] | #xb7 | [#x300-#x36f] | [#x203f-#x2040]
// PN_CHARS_BASE ::= [A-Z] | [a-z] | [#xc0-#xd6] | [#xd8-#xf6] | [#xf8-#x2ff] | [#x370-#x37d] | [#x37f-#x1fff] | [#x200c-#x200d] | [#x2070-#x218f] | [#x2c00-#x2fef] | [#x3001-#xd7ff] | [#xf900-#xfdcf] | [#xfdf0-#xfffd] | [#x10000-#xeffff]
// PN_CHARS_U ::= PN_CHARS_BASE | '_'
func IsRune_PN_CHARS(r rune) bool {
	switch {
	case r == '-':
		return true
	case '0' <= r && r <= '9':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	case r == '_':
		return true
	case 'a' <= r && r <= 'z':
		return true
	case r == 0xB7:
		return true
	case 0xC0 <= r && r <= 0xD6:
		return true
	case 0xD8 <= r && r <= 0xF6:
		return true
	case 0xF8 <= r && r <= 0x37D:
		return true
	case 0x37F <= r && r <= 0x1FFF:
		return true
	case 0x200C <= r && r <= 0x200D:
		return true
	case 0x203F <= r && r <= 0x2040:
		return true
	case 0x2070 <= r && r <= 0x218F:
		return true
	case 0x2C00 <= r && r <= 0x2FEF:
		return true
	case 0x3001 <= r && r <= 0xD7FF:
		return true
	case 0xF900 <= r && r <= 0xFDCF:
		return true
	case 0xFDF0 <= r && r <= 0xFFFD:
		return true
	case 0x10000 <= r && r <= 0xEFFFF:
		return true
	}

	return false
}
