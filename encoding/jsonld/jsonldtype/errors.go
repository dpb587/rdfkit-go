package jsonldtype

import "fmt"

type Error struct {
	Code ErrorCode
	Err  error
}

var _ error = Error{}

func (e Error) Error() string {
	if e.Err == nil {
		return string(e.Code)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Err.Error())
}

func (e Error) Unwrap() error {
	return e.Err
}

//

type ErrorCode string

const (
	CollidingKeywords           ErrorCode = "colliding keywords"
	ConflictingIndexes          ErrorCode = "conflicting indexes"
	ContextOverflow             ErrorCode = "context overflow"
	CyclicIRIMapping            ErrorCode = "cyclic IRI mapping"
	InvalidAtIDValue            ErrorCode = "invalid @id value"
	InvalidAtImportValue        ErrorCode = "invalid @import value"
	InvalidAtIncludedValue      ErrorCode = "invalid @included value"
	InvalidAtIndexValue         ErrorCode = "invalid @index value"
	InvalidAtNestValue          ErrorCode = "invalid @nest value"
	InvalidAtPrefixValue        ErrorCode = "invalid @prefix value"
	InvalidAtPropagateValue     ErrorCode = "invalid @propagate value"
	InvalidAtProtectedValue     ErrorCode = "invalid @protected value"
	InvalidAtReverseValue       ErrorCode = "invalid @reverse value"
	InvalidAtVersionValue       ErrorCode = "invalid @version value"
	InvalidBaseDirection        ErrorCode = "invalid base direction"
	InvalidBaseIRI              ErrorCode = "invalid base IRI"
	InvalidContainerMapping     ErrorCode = "invalid container mapping"
	InvalidContextEntry         ErrorCode = "invalid context entry"
	InvalidContextNullification ErrorCode = "invalid context nullification"
	InvalidDefaultLanguage      ErrorCode = "invalid default language"
	InvalidIRIMapping           ErrorCode = "invalid IRI mapping"
	InvalidJSONLiteral          ErrorCode = "invalid JSON literal"
	InvalidKeywordAlias         ErrorCode = "invalid keyword alias"
	InvalidLanguageMapValue     ErrorCode = "invalid language map value"
	InvalidLanguageMapping      ErrorCode = "invalid language mapping"
	InvalidLanguageTaggedString ErrorCode = "invalid language-tagged string"
	InvalidLanguageTaggedValue  ErrorCode = "invalid language-tagged value"
	InvalidLocalContext         ErrorCode = "invalid local context"
	InvalidRemoteContext        ErrorCode = "invalid remote context"
	InvalidReversePropertyMap   ErrorCode = "invalid reverse property map"
	InvalidReversePropertyValue ErrorCode = "invalid reverse property value"
	InvalidReverseProperty      ErrorCode = "invalid reverse property"
	InvalidScopedContext        ErrorCode = "invalid scoped context"
	InvalidScriptElement        ErrorCode = "invalid script element"
	InvalidSetOrListObject      ErrorCode = "invalid set or list object"
	InvalidTermDefinition       ErrorCode = "invalid term definition"
	InvalidTypeMapping          ErrorCode = "invalid type mapping"
	InvalidTypeValue            ErrorCode = "invalid type value"
	InvalidTypedValue           ErrorCode = "invalid typed value"
	InvalidValueObjectValue     ErrorCode = "invalid value object value"
	InvalidValueObject          ErrorCode = "invalid value object"
	InvalidVocabMapping         ErrorCode = "invalid vocab mapping"
	IRIConfusedWithPrefix       ErrorCode = "IRI confused with prefix"
	KeywordRedefinition         ErrorCode = "keyword redefinition"
	LoadingDocumentFailed       ErrorCode = "loading document failed"
	LoadingRemoteContextFailed  ErrorCode = "loading remote context failed"
	MultipleContextLinkHeaders  ErrorCode = "multiple context link headers"
	ProcessingModeConflict      ErrorCode = "processing mode conflict"
	ProtectedTermRedefinition   ErrorCode = "protected term redefinition"
)
