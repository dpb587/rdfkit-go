package xsdobject

import (
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdtype"
	"github.com/dpb587/rdfkit-go/rdf"
)

func AnyURI(v string) rdf.ObjectValue {
	return xsdtype.AnyURI(v).AsObjectValue()
}

func Base64Binary(v []byte) rdf.ObjectValue {
	return xsdtype.Base64Binary(v).AsObjectValue()
}

func Boolean(v bool) rdf.ObjectValue {
	return xsdtype.Boolean(v).AsObjectValue()
}

func Byte(v byte) rdf.ObjectValue {
	return xsdtype.Byte(v).AsObjectValue()
}

func DateTimeStamp(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.DateTimeStamp{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

func DateTime(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.DateTime{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

func Date(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.Date{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

// TODO Duration

func Decimal(v float64) rdf.ObjectValue {
	return xsdtype.Decimal(v).AsObjectValue()
}

func Double(v float64) rdf.ObjectValue {
	return xsdtype.Double(v).AsObjectValue()
}

func Float(v float32) rdf.ObjectValue {
	return xsdtype.Float(v).AsObjectValue()
}

func GDay(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.GDay{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

func GMonthDay(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.GMonthDay{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

func GMonth(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.GMonth{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

func GYearMonth(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.GYearMonth{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

func GYear(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.GYear{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

func HexBinary(v []byte) rdf.ObjectValue {
	return xsdtype.HexBinary(v).AsObjectValue()
}

func Int(v int32) rdf.ObjectValue {
	return xsdtype.Int(v).AsObjectValue()
}

func Integer(v int64) rdf.ObjectValue {
	return xsdtype.Integer(v).AsObjectValue()
}

func Long(v int64) rdf.ObjectValue {
	return xsdtype.Long(v).AsObjectValue()
}

func Short(v int16) rdf.ObjectValue {
	return xsdtype.Short(v).AsObjectValue()
}

func String(v string) rdf.ObjectValue {
	return xsdtype.String(v).AsObjectValue()
}

func Time(layout string, value time.Time) rdf.ObjectValue {
	return xsdtype.Time{
		Time:   value,
		Layout: layout,
	}.AsObjectValue()
}

func UnsignedByte(v uint8) rdf.ObjectValue {
	return xsdtype.UnsignedByte(v).AsObjectValue()
}

func UnsignedInt(v uint32) rdf.ObjectValue {
	return xsdtype.UnsignedInt(v).AsObjectValue()
}

func UnsignedLong(v uint64) rdf.ObjectValue {
	return xsdtype.UnsignedLong(v).AsObjectValue()
}

func UnsignedShort(v uint16) rdf.ObjectValue {
	return xsdtype.UnsignedShort(v).AsObjectValue()
}
