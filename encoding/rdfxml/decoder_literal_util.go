package rdfxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

func (d *Decoder) xmlRender() ([]byte, xml.EndElement, error) {
	buf := bytes.NewBuffer(nil)

	// TODO xmlDecoder.Namespaces()

	xmlEncoder := xml.NewEncoder(buf)

	var elementDepth int

	for {
		rawToken, err := d.tokenNext()
		if err != nil {
			return nil, xml.EndElement{}, fmt.Errorf("read token: %v", err)
		}

		switch tT := rawToken.(type) {
		case xml.StartElement:
			elementDepth++
		case xml.EndElement:
			if elementDepth == 0 {
				err := xmlEncoder.Flush()
				if err != nil {
					return nil, xml.EndElement{}, fmt.Errorf("flush: %v", err)
				}

				return buf.Bytes(), tT, nil
			}

			elementDepth--
		}

		err = xmlEncoder.EncodeToken(rawToken)
		if err != nil {
			return nil, xml.EndElement{}, fmt.Errorf("write token: %v", err)
		}
	}
}
