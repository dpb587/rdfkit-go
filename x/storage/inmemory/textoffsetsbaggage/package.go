package textoffsetsbaggage

import (
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/x/storage/inmemory"
)

type baggageKeyType struct{}

var baggageKey = baggageKeyType{}

func Add(s *inmemory.Statement, offsets encoding.StatementTextOffsets) {
	if len(offsets) == 0 {
		return
	} else if s.Baggage == nil {
		s.Baggage = map[any]any{
			baggageKey: []encoding.StatementTextOffsets{
				offsets,
			},
		}

		return
	}

	existing, _ := s.Baggage[baggageKey].([]encoding.StatementTextOffsets)
	s.Baggage[baggageKey] = append(existing, offsets)
}

func AddAll(s *inmemory.Statement, offsets []encoding.StatementTextOffsets) {
	for _, o := range offsets {
		Add(s, o)
	}
}

func Get(s *inmemory.Statement) []encoding.StatementTextOffsets {
	if s.Baggage == nil {
		return nil
	}

	existing, _ := s.Baggage[baggageKey].([]encoding.StatementTextOffsets)
	return existing
}
