package internal

const (
	HexUpper = "0123456789ABCDEF"
	HexLower = "0123456789abcdef"
)

func HexDecode(c rune) (rune, bool) {
	switch c {
	case '0':
		return 0, true
	case '1':
		return 1, true
	case '2':
		return 2, true
	case '3':
		return 3, true
	case '4':
		return 4, true
	case '5':
		return 5, true
	case '6':
		return 6, true
	case '7':
		return 7, true
	case '8':
		return 8, true
	case '9':
		return 9, true
	case 'A', 'a':
		return 10, true
	case 'B', 'b':
		return 11, true
	case 'C', 'c':
		return 12, true
	case 'D', 'd':
		return 13, true
	case 'E', 'e':
		return 14, true
	case 'F', 'f':
		return 15, true
	}

	return 0, false
}
