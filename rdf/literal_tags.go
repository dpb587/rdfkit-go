package rdf

// LiteralTag represents additional metadata about a [Literal] term. They should only be used in specific contexts,
// typically dependent on the literal's datatype.
//
// This is a closed interface. See [LanguageLiteralTag] and [DirectionalLanguageLiteralTag].
type LiteralTag interface {
	Equals(other LiteralTag) bool

	isLiteralTag()
}

//

// LanguageLiteralTag is used for a language-tagged literal.
type LanguageLiteralTag struct {
	// Language is a well-formed [BCP47] language tag.
	//
	// [BCP47]: https://www.rfc-editor.org/rfc/rfc5646
	Language string
}

var _ LiteralTag = LanguageLiteralTag{}

func (LanguageLiteralTag) isLiteralTag() {}

func (t LanguageLiteralTag) Equals(other LiteralTag) bool {
	otherTag, ok := other.(LanguageLiteralTag)
	if !ok {
		return false
	}

	return t.Language == otherTag.Language
}

//

// DirectionalLanguageLiteralTag behavior is still undefined for this Go project and pending RDF 1.2. It is currently
// used internally for testing and potential future use.
type DirectionalLanguageLiteralTag struct {
	// Language is a well-formed [BCP47] language tag.
	//
	// [BCP47]: https://www.rfc-editor.org/rfc/rfc5646
	Language string

	// BaseDirection indicates the initial text direction.
	//
	// Valid values are "ltr" and "rtl".
	BaseDirection string
}

var _ LiteralTag = DirectionalLanguageLiteralTag{}

func (DirectionalLanguageLiteralTag) isLiteralTag() {}

func (t DirectionalLanguageLiteralTag) Equals(other LiteralTag) bool {
	otherTag, ok := other.(DirectionalLanguageLiteralTag)
	if !ok {
		return false
	}

	return t.Language == otherTag.Language && t.BaseDirection == otherTag.BaseDirection
}
