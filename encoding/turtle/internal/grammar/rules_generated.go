// Code generated by rulegencmd; DO NOT EDIT.

package grammar

import "fmt"

type R uint

const (
	// R_turtleDoc ::= statement*
	R_turtleDoc R = iota

	// R_statement ::= directive | ( triples L_PERIOD )
	R_statement R = iota

	// R_directive ::= prefixID | base | sparqlPrefix | sparqlBase
	R_directive R = iota

	// R_prefixID ::= '@prefix' PNAME_NS IRIREF L_PERIOD
	R_prefixID R = iota

	// R_base ::= '@base' IRIREF L_PERIOD
	R_base R = iota

	// R_sparqlBase ::= ~'BASE' IRIREF
	R_sparqlBase R = iota

	// R_sparqlPrefix ::= ~'PREFIX' PNAME_NS IRIREF
	R_sparqlPrefix R = iota

	// R_triples ::= ( subject predicateObjectList ) | ( blankNodePropertyList predicateObjectList? )
	R_triples R = iota

	// R_predicateObjectList ::= verb objectList ( ( L_SEMICOLON ( ( verb objectList )? ) )* )
	R_predicateObjectList R = iota

	// R_objectList ::= object ( ( L_COMMA object )* )
	R_objectList R = iota

	// R_verb ::= predicate | 'a'
	R_verb R = iota

	// R_subject ::= iri | BlankNode | collection
	R_subject R = iota

	// R_predicate ::= iri
	R_predicate R = iota

	// R_object ::= iri | BlankNode | collection | blankNodePropertyList | literal
	R_object R = iota

	// R_literal ::= RDFLiteral | NumericLiteral | BooleanLiteral
	R_literal R = iota

	// R_blankNodePropertyList ::= L_OPEN_BRACKET predicateObjectList L_CLOSE_BRACKET
	R_blankNodePropertyList R = iota

	// R_collection ::= L_OPEN_PAREN object* L_CLOSE_PAREN
	R_collection R = iota

	// R_NumericLiteral ::= INTEGER | DECIMAL | DOUBLE
	R_NumericLiteral R = iota

	// R_RDFLiteral ::= String ( ( LANGTAG | ( L_DOUBLE_CARET iri ) )? )
	R_RDFLiteral R = iota

	// R_BooleanLiteral ::= 'true' | 'false'
	R_BooleanLiteral R = iota

	// R_String ::= STRING_LITERAL_QUOTE | STRING_LITERAL_SINGLE_QUOTE | STRING_LITERAL_LONG_SINGLE_QUOTE | STRING_LITERAL_LONG_QUOTE
	R_String R = iota

	// R_iri ::= IRIREF | PrefixedName
	R_iri R = iota

	// R_PrefixedName ::= PNAME_LN | PNAME_NS
	R_PrefixedName R = iota

	// R_BlankNode ::= BLANK_NODE_LABEL | ANON
	R_BlankNode R = iota

	// R_IRIREF ::= '<' ( ( [^#x0-#x20<>"{}|^`\] | UCHAR )* ) '>'
	R_IRIREF R = iota

	// R_PNAME_NS ::= PN_PREFIX? ':'
	R_PNAME_NS R = iota

	// R_PNAME_LN ::= PNAME_NS PN_LOCAL
	R_PNAME_LN R = iota

	// R_BLANK_NODE_LABEL ::= '_:' ( PN_CHARS_U | [0-9] ) ( ( ( ( PN_CHARS | L_PERIOD )* ) PN_CHARS )? )
	R_BLANK_NODE_LABEL R = iota

	// R_LANGTAG ::= '@' ( ( [a-z] | [A-Z] )+ ) ( ( '-' ( ( [a-z] | [A-Z] | [0-9] )+ ) )* )
	R_LANGTAG R = iota

	// R_INTEGER ::= [+-]? [0-9]+
	R_INTEGER R = iota

	// R_DECIMAL ::= [+-]? [0-9]* '.' [0-9]+
	R_DECIMAL R = iota

	// R_DOUBLE ::= [+-]? ( ( [0-9]+ L_PERIOD [0-9]* EXPONENT ) | ( L_PERIOD [0-9]+ EXPONENT ) | ( [0-9]+ EXPONENT ) )
	R_DOUBLE R = iota

	// R_EXPONENT ::= [eE] [+-]? [0-9]+
	R_EXPONENT R = iota

	// R_STRING_LITERAL_QUOTE ::= '"' ( ( [^"\#xa#xd] | ECHAR | UCHAR )* ) '"'
	R_STRING_LITERAL_QUOTE R = iota

	// R_STRING_LITERAL_SINGLE_QUOTE ::= "'" ( ( [^'\#xa#xd] | ECHAR | UCHAR )* ) "'"
	R_STRING_LITERAL_SINGLE_QUOTE R = iota

	// R_STRING_LITERAL_LONG_SINGLE_QUOTE ::= "'''" ( ( ( ( "'" | "''" )? ) ( [^'\\] | ECHAR | UCHAR ) )* ) "'''"
	R_STRING_LITERAL_LONG_SINGLE_QUOTE R = iota

	// R_STRING_LITERAL_LONG_QUOTE ::= '"""' ( ( ( ( '"' | '""' )? ) ( [^"\\] | ECHAR | UCHAR ) )* ) '"""'
	R_STRING_LITERAL_LONG_QUOTE R = iota

	// R_UCHAR ::= ( '\u' HEX HEX HEX HEX ) | ( '\U' HEX HEX HEX HEX HEX HEX HEX HEX )
	R_UCHAR R = iota

	// R_ECHAR ::= '\' [tbnrf"'\]
	R_ECHAR R = iota

	// R_WS ::= ' ' | #x9 | #xd | #xa
	R_WS R = iota

	// R_ANON ::= '[' WS* ']'
	R_ANON R = iota

	// R_PN_CHARS_BASE ::= [A-Z] | [a-z] | [#xc0-#xd6] | [#xd8-#xf6] | [#xf8-#x2ff] | [#x370-#x37d] | [#x37f-#x1fff] | [#x200c-#x200d] | [#x2070-#x218f] | [#x2c00-#x2fef] | [#x3001-#xd7ff] | [#xf900-#xfdcf] | [#xfdf0-#xfffd] | [#x10000-#xeffff]
	R_PN_CHARS_BASE R = iota

	// R_PN_CHARS_U ::= PN_CHARS_BASE | '_'
	R_PN_CHARS_U R = iota

	// R_PN_CHARS ::= PN_CHARS_U | '-' | [0-9] | #xb7 | [#x300-#x36f] | [#x203f-#x2040]
	R_PN_CHARS R = iota

	// R_PN_PREFIX ::= PN_CHARS_BASE ( ( ( ( PN_CHARS | '.' )* ) PN_CHARS )? )
	R_PN_PREFIX R = iota

	// R_PN_LOCAL ::= ( PN_CHARS_U | ':' | [0-9] | PLX ) ( ( ( ( PN_CHARS | '.' | ':' | PLX )* ) ( PN_CHARS | ':' | PLX ) )? )
	R_PN_LOCAL R = iota

	// R_PLX ::= PERCENT | PN_LOCAL_ESC
	R_PLX R = iota

	// R_PERCENT ::= '%' HEX HEX
	R_PERCENT R = iota

	// R_HEX ::= [0-9] | [A-F] | [a-f]
	R_HEX R = iota

	// R_PN_LOCAL_ESC ::= '\' ( '_' | '~' | '.' | '-' | '!' | '$' | '&' | "'" | '(' | ')' | '*' | '+' | ',' | ';' | '=' | '/' | '?' | '#' | '@' | '%' )
	R_PN_LOCAL_ESC R = iota

	// R_WS_COMMENT ::= '#' [^#xd#xa]*
	R_WS_COMMENT R = iota

	// R_L_PERIOD ::= '.'
	R_L_PERIOD R = iota

	// R_L_SEMICOLON ::= ';'
	R_L_SEMICOLON R = iota

	// R_L_COMMA ::= ','
	R_L_COMMA R = iota

	// R_L_OPEN_BRACKET ::= '['
	R_L_OPEN_BRACKET R = iota

	// R_L_CLOSE_BRACKET ::= ']'
	R_L_CLOSE_BRACKET R = iota

	// R_L_OPEN_PAREN ::= '('
	R_L_OPEN_PAREN R = iota

	// R_L_CLOSE_PAREN ::= ')'
	R_L_CLOSE_PAREN R = iota

	// R_L_OPEN_BRACE ::= '{'
	R_L_OPEN_BRACE R = iota

	// R_L_CLOSE_BRACE ::= '}'
	R_L_CLOSE_BRACE R = iota

	// R_L_DOUBLE_CARET ::= '^^'
	R_L_DOUBLE_CARET R = iota

	// R_requiredWS ::= ( WS | WS_COMMENT )+
	R_requiredWS R = iota

	// R_optionalWS ::= ( WS | WS_COMMENT )*
	R_optionalWS R = iota
)

func (r R) String() string {
	switch r {
	case R_turtleDoc:
		return "turtleDoc"
	case R_statement:
		return "statement"
	case R_directive:
		return "directive"
	case R_prefixID:
		return "prefixID"
	case R_base:
		return "base"
	case R_sparqlBase:
		return "sparqlBase"
	case R_sparqlPrefix:
		return "sparqlPrefix"
	case R_triples:
		return "triples"
	case R_predicateObjectList:
		return "predicateObjectList"
	case R_objectList:
		return "objectList"
	case R_verb:
		return "verb"
	case R_subject:
		return "subject"
	case R_predicate:
		return "predicate"
	case R_object:
		return "object"
	case R_literal:
		return "literal"
	case R_blankNodePropertyList:
		return "blankNodePropertyList"
	case R_collection:
		return "collection"
	case R_NumericLiteral:
		return "NumericLiteral"
	case R_RDFLiteral:
		return "RDFLiteral"
	case R_BooleanLiteral:
		return "BooleanLiteral"
	case R_String:
		return "String"
	case R_iri:
		return "iri"
	case R_PrefixedName:
		return "PrefixedName"
	case R_BlankNode:
		return "BlankNode"
	case R_IRIREF:
		return "IRIREF"
	case R_PNAME_NS:
		return "PNAME_NS"
	case R_PNAME_LN:
		return "PNAME_LN"
	case R_BLANK_NODE_LABEL:
		return "BLANK_NODE_LABEL"
	case R_LANGTAG:
		return "LANGTAG"
	case R_INTEGER:
		return "INTEGER"
	case R_DECIMAL:
		return "DECIMAL"
	case R_DOUBLE:
		return "DOUBLE"
	case R_EXPONENT:
		return "EXPONENT"
	case R_STRING_LITERAL_QUOTE:
		return "STRING_LITERAL_QUOTE"
	case R_STRING_LITERAL_SINGLE_QUOTE:
		return "STRING_LITERAL_SINGLE_QUOTE"
	case R_STRING_LITERAL_LONG_SINGLE_QUOTE:
		return "STRING_LITERAL_LONG_SINGLE_QUOTE"
	case R_STRING_LITERAL_LONG_QUOTE:
		return "STRING_LITERAL_LONG_QUOTE"
	case R_UCHAR:
		return "UCHAR"
	case R_ECHAR:
		return "ECHAR"
	case R_WS:
		return "WS"
	case R_ANON:
		return "ANON"
	case R_PN_CHARS_BASE:
		return "PN_CHARS_BASE"
	case R_PN_CHARS_U:
		return "PN_CHARS_U"
	case R_PN_CHARS:
		return "PN_CHARS"
	case R_PN_PREFIX:
		return "PN_PREFIX"
	case R_PN_LOCAL:
		return "PN_LOCAL"
	case R_PLX:
		return "PLX"
	case R_PERCENT:
		return "PERCENT"
	case R_HEX:
		return "HEX"
	case R_PN_LOCAL_ESC:
		return "PN_LOCAL_ESC"
	case R_WS_COMMENT:
		return "WS_COMMENT"
	case R_L_PERIOD:
		return "L_PERIOD"
	case R_L_SEMICOLON:
		return "L_SEMICOLON"
	case R_L_COMMA:
		return "L_COMMA"
	case R_L_OPEN_BRACKET:
		return "L_OPEN_BRACKET"
	case R_L_CLOSE_BRACKET:
		return "L_CLOSE_BRACKET"
	case R_L_OPEN_PAREN:
		return "L_OPEN_PAREN"
	case R_L_CLOSE_PAREN:
		return "L_CLOSE_PAREN"
	case R_L_OPEN_BRACE:
		return "L_OPEN_BRACE"
	case R_L_CLOSE_BRACE:
		return "L_CLOSE_BRACE"
	case R_L_DOUBLE_CARET:
		return "L_DOUBLE_CARET"
	case R_requiredWS:
		return "requiredWS"
	case R_optionalWS:
		return "optionalWS"
	}

	return fmt.Sprintf("grammar.R(%d)", r)
}
