package schemaobject

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/schema/schematype"
	"github.com/dpb587/rdfkit-go/rdf"
)

func Boolean(v bool) rdf.ObjectValue {
	return schematype.Boolean(v).AsObjectValue()
}

func CssSelectorType(v string) rdf.ObjectValue {
	return schematype.CssSelectorType(v).AsObjectValue()
}

func DateTime(layout string, time time.Time) rdf.ObjectValue {
	return schematype.DateTime{
		Time:   time,
		Layout: layout,
	}.AsObjectValue()
}

func Date(layout string, time time.Time) rdf.ObjectValue {
	return schematype.Date{
		Time:   time,
		Layout: layout,
	}.AsObjectValue()
}

func Float(v float32) rdf.ObjectValue {
	return schematype.Float(v).AsObjectValue()
}

func Integer(v int64) rdf.ObjectValue {
	return schematype.Integer(v).AsObjectValue()
}

func Number(v float64) rdf.ObjectValue {
	return schematype.Number(v).AsObjectValue()
}

func Text(v string) rdf.ObjectValue {
	return schematype.Text(v).AsObjectValue()
}

func URL(v string) rdf.ObjectValue {
	return schematype.URL(v).AsObjectValue()
}

func XPathType(v string) rdf.ObjectValue {
	return schematype.XPathType(v).AsObjectValue()
}
