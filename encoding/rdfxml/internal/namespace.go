package internal

// https://www.w3.org/TR/rdf-syntax-grammar/#section-Namespace

const Space = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"

const (
	Local_RDF_Syntax         = "RDF"
	Local_Description_Syntax = "Description"
	Local_ID_Syntax          = "ID"
	Local_About_Syntax       = "about"
	Local_ParseType_Syntax   = "parseType"
	Local_Resource_Syntax    = "resource"
	Local_Li_Syntax          = "li"
	Local_NodeID_Syntax      = "nodeID"
	Local_Datatype_Syntax    = "datatype"

	Local_Seq_Class        = "Seq"
	Local_Bag_Class        = "Bag"
	Local_Alt_Class        = "Alt"
	Local_Statement_Classs = "Statement"
	Local_Property_Class   = "Property"
	Local_XMLLiteral_Class = "XMLLiteral"
	Local_List_Class       = "List"

	Local_Subject_Property   = "subject"
	Local_Predicate_Property = "predicate"
	Local_Object_Property    = "object"
	Local_Type_Property      = "type"
	Local_Value_Property     = "value"
	Local_First_Property     = "first"
	Local_Rest_Property      = "rest"
	// _n_Property = "_n" // where n is a decimal integer greater than zero with no leading zeros

	Local_Nil_Resource = "nil"

	// withdrawn
	Local_AboutEach_Old       = "aboutEach"
	Local_AboutEachPrefix_Old = "aboutEachPrefix"
	Local_BagID_Old           = "bagID"
)
