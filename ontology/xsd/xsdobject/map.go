package xsdobject

import (
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdtype"
	"github.com/dpb587/rdfkit-go/rdf"
)

func MapAnyURI(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapAnyURI(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapBase64Binary(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapBase64Binary(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapBoolean(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapBoolean(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapByte(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapByte(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapDateTimeStamp(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapDateTimeStamp(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapDateTime(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapDateTime(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapDate(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapDate(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapDuration(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapDuration(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapDecimal(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapDecimal(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapDouble(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapDouble(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapFloat(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapFloat(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapGDay(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapGDay(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapGMonthDay(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapGMonthDay(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapGMonth(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapGMonth(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapGYearMonth(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapGYearMonth(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapGYear(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapGYear(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapHexBinary(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapHexBinary(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapInt(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapInt(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapInteger(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapInteger(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapLong(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapLong(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapShort(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapShort(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapString(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapString(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapTime(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapTime(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapUnsignedByte(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapUnsignedByte(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapUnsignedInt(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapUnsignedInt(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapUnsignedLong(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapUnsignedLong(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}

func MapUnsignedShort(lexicalForm string) (rdf.ObjectValue, error) {
	v, err := xsdtype.MapUnsignedShort(lexicalForm)
	if err != nil {
		return nil, err
	}

	return v.AsObjectValue(), nil
}
