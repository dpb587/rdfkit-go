package jsonldtype

type UnknownKeywordError struct {
	Keyword string
}

func (e *UnknownKeywordError) Error() string {
	return "unknown keyword: " + e.Keyword
}
