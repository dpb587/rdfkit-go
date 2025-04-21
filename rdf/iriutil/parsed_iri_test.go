package iriutil

import (
	"net/url"
	"testing"
)

func TestParsedIRI_PreserveRawPath(t *testing.T) {
	original := "http://example.com/Dürst"

	{ // std; escapes ü
		vv, err := url.Parse(original)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if _a, _e := vv.String(), "http://example.com/D%C3%BCrst"; _a != _e {
			t.Fatalf("url string: expected %v, got %v", _e, _a)
		}
	}

	vv, err := ParseIRI(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _a, _e := vv.String(), original; _a != _e {
		t.Fatalf("iri string: expected %v, got %v", _e, _a)
	}
}

func TestParsedIRI_PreserveRawFragment(t *testing.T) {
	original := "http://www.w3.org/2013/RDFXMLTests/rdfms-difference-between-ID-and-about/test3.rdf#Dürst"

	{ // std; escapes ü
		vv, err := url.Parse(original)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if _a, _e := vv.String(), "http://www.w3.org/2013/RDFXMLTests/rdfms-difference-between-ID-and-about/test3.rdf#D%C3%BCrst"; _a != _e {
			t.Fatalf("url string: expected %v, got %v", _e, _a)
		}
	}

	vv, err := ParseIRI(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _a, _e := vv.String(), original; _a != _e {
		t.Fatalf("iri string: expected %v, got %v", _e, _a)
	}
}

func TestParsedIRI_PreserveEmptyFragment(t *testing.T) {
	original := "http://example.com/#"

	{ // std; drops #
		vv, err := url.Parse(original)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if _a, _e := vv.String(), "http://example.com/"; _a != _e {
			t.Fatalf("url string: expected %v, got %v", _e, _a)
		}
	}

	vv, err := ParseIRI(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _a, _e := vv.String(), original; _a != _e {
		t.Fatalf("iri string: expected %v, got %v", _e, _a)
	}
}

func TestParsedIRI_ParsePreserveRawPath(t *testing.T) {
	base := "http://example.com/"

	{ // std; escapes ü
		vv, err := url.Parse(base)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		vv, err = vv.Parse("Dürst")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if _a, _e := vv.String(), "http://example.com/D%C3%BCrst"; _a != _e {
			t.Fatalf("url string: expected %v, got %v", _e, _a)
		}
	}

	vv, err := ParseIRI(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vv, err = vv.Parse("Dürst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _a, _e := vv.String(), "http://example.com/Dürst"; _a != _e {
		t.Fatalf("iri string: expected %v, got %v", _e, _a)
	}
}

func TestParsedIRI_ParsePreserveRawFragment(t *testing.T) {
	base := "http://example.com/"

	{ // std; escapes ü
		vv, err := url.Parse(base)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		vv, err = vv.Parse("#Dürst")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if _a, _e := vv.String(), "http://example.com/#D%C3%BCrst"; _a != _e {
			t.Fatalf("url string: expected %v, got %v", _e, _a)
		}
	}

	vv, err := ParseIRI(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vv, err = vv.Parse("#Dürst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _a, _e := vv.String(), "http://example.com/#Dürst"; _a != _e {
		t.Fatalf("iri string: expected %v, got %v", _e, _a)
	}
}

func TestParsedIRI_CoverageClearFragment(t *testing.T) {
	base := "http://example.com/#anything"

	{
		vv, err := url.Parse(base)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		vv, err = vv.Parse("path")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if _a, _e := vv.String(), "http://example.com/path"; _a != _e {
			t.Fatalf("url string: expected %v, got %v", _e, _a)
		}
	}

	vv, err := ParseIRI(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vv, err = vv.Parse("path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _a, _e := vv.String(), "http://example.com/path"; _a != _e {
		t.Fatalf("iri string: expected %v, got %v", _e, _a)
	}
}

func TestParsedIRI_CoverageMerge(t *testing.T) {
	base := "http://example.com/#anything"

	{ // std; escapes ü
		vv, err := url.Parse(base)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		vv, err = vv.Parse("Dürst#Dürst")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if _a, _e := vv.String(), "http://example.com/D%C3%BCrst#D%C3%BCrst"; _a != _e {
			t.Fatalf("url string: expected %v, got %v", _e, _a)
		}
	}

	vv, err := ParseIRI(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vv, err = vv.Parse("Dürst#Dürst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _a, _e := vv.String(), "http://example.com/Dürst#Dürst"; _a != _e {
		t.Fatalf("iri string: expected %v, got %v", _e, _a)
	}
}

func TestParsedIRI_CoverageParseAbs(t *testing.T) {
	base := "http://example.com/#anything"

	{
		vv, err := url.Parse(base)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		vv, err = vv.Parse("http://two.example.com/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		} else if _a, _e := vv.String(), "http://two.example.com/"; _a != _e {
			t.Fatalf("url string: expected %v, got %v", _e, _a)
		}
	}

	vv, err := ParseIRI(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vv, err = vv.Parse("http://two.example.com/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _a, _e := vv.String(), "http://two.example.com/"; _a != _e {
		t.Fatalf("iri string: expected %v, got %v", _e, _a)
	}
}
