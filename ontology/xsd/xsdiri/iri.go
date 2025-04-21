package xsdiri

import "github.com/dpb587/rdfkit-go/rdf"

// See: https://www.w3.org/TR/xmlschema11-2/#built-in-datatypes

const (
	Base rdf.IRI = "http://www.w3.org/2001/XMLSchema#"

	AnyURI_Datatype             = Base + "anyURI"             // VALUE
	Base64Binary_Datatype       = Base + "base64Binary"       // VALUE
	Boolean_Datatype            = Base + "boolean"            // VALUE
	Byte_Datatype               = Base + "byte"               // VALUE
	Date_Datatype               = Base + "date"               // VALUE
	DateTime_Datatype           = Base + "dateTime"           // VALUE
	DateTimeStamp_Datatype      = Base + "dateTimeStamp"      // VALUE
	DayTimeDuration_Datatype    = Base + "dayTimeDuration"    //
	Decimal_Datatype            = Base + "decimal"            // VALUE
	Double_Datatype             = Base + "double"             // VALUE
	Duration_Datatype           = Base + "duration"           //
	ENTITY_Datatype             = Base + "ENTITY"             //
	Float_Datatype              = Base + "float"              //
	GDay_Datatype               = Base + "gDay"               // VALUE
	GMonth_Datatype             = Base + "gMonth"             // VALUE
	GMonthDay_Datatype          = Base + "gMonthDay"          // VALUE
	GYear_Datatype              = Base + "gYear"              // VALUE
	GYearMonth_Datatype         = Base + "gYearMonth"         // VALUE
	HexBinary_Datatype          = Base + "hexBinary"          //
	ID_Datatype                 = Base + "ID"                 //
	IDREF_Datatype              = Base + "IDREF"              //
	Int_Datatype                = Base + "int"                //
	Integer_Datatype            = Base + "integer"            // VALUE
	Language_Datatype           = Base + "language"           //
	Long_Datatype               = Base + "long"               // VALUE
	Name_Datatype               = Base + "Name"               //
	NCName_Datatype             = Base + "NCName"             //
	NegativeInteger_Datatype    = Base + "negativeInteger"    //
	NMTOKEN_Datatype            = Base + "NMTOKEN"            //
	NonNegativeInteger_Datatype = Base + "nonNegativeInteger" //
	NonPositiveInteger_Datatype = Base + "nonPositiveInteger" //
	NormalizedString_Datatype   = Base + "normalizedString"   //
	NOTATION_Datatype           = Base + "NOTATION"           //
	PositiveInteger_Datatype    = Base + "positiveInteger"    //
	QName_Datatype              = Base + "QName"              //
	Short_Datatype              = Base + "short"              // VALUE
	String_Datatype             = Base + "string"             // VALUE
	Time_Datatype               = Base + "time"               // VALUE
	Token_Datatype              = Base + "token"              //
	UnsignedByte_Datatype       = Base + "unsignedByte"       //
	UnsignedInt_Datatype        = Base + "unsignedInt"        //
	UnsignedLong_Datatype       = Base + "unsignedLong"       //
	UnsignedShort_Datatype      = Base + "unsignedShort"      //
	YearMonthDuration_Datatype  = Base + "yearMonthDuration"  //
)
