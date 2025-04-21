package jsonldinternal

var definedKeywords = map[string]struct{}{
	"@base":      {},
	"@container": {},
	"@context":   {},
	"@direction": {},
	"@graph":     {},
	"@id":        {},
	"@import":    {},
	"@included":  {},
	"@index":     {},
	"@json":      {},
	"@language":  {},
	"@list":      {},
	"@nest":      {},
	"@none":      {},
	"@prefix":    {},
	"@propagate": {},
	"@protected": {},
	"@reverse":   {},
	"@set":       {},
	"@type":      {},
	"@value":     {},
	"@version":   {},
	"@vocab":     {},
}

func IsKeyword(k string) bool {
	_, ok := definedKeywords[k]
	return ok
}
