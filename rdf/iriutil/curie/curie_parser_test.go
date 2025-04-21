package curie

import "testing"

func TestParseCURIE(t *testing.T) {
	for _, tc := range []struct {
		Input         string
		ExpectedCURIE CURIE
		ExpectedOK    bool
	}{
		{
			Input:      "",
			ExpectedOK: false,
		},
		{
			Input: "ex:foo",
			ExpectedCURIE: CURIE{
				Prefix:    "ex",
				Reference: "foo",
			},
			ExpectedOK: true,
		},
		{
			Input: ":foo",
			ExpectedCURIE: CURIE{
				Prefix:    "",
				Reference: "foo",
			},
			ExpectedOK: true,
		},
		{
			Input: "foo",
			ExpectedCURIE: CURIE{
				DefaultPrefix: true,
				Reference:     "foo",
			},
			ExpectedOK: true,
		},
		{
			Input: "[ex:foo]",
			ExpectedCURIE: CURIE{
				Safe:      true,
				Prefix:    "ex",
				Reference: "foo",
			},
			ExpectedOK: true,
		},
		{
			Input: "[:foo]",
			ExpectedCURIE: CURIE{
				Safe:      true,
				Prefix:    "",
				Reference: "foo",
			},
			ExpectedOK: true,
		},
		{
			Input: "[foo]",
			ExpectedCURIE: CURIE{
				Safe:          true,
				DefaultPrefix: true,
				Reference:     "foo",
			},
			ExpectedOK: true,
		},
	} {
		t.Run(tc.Input, func(t *testing.T) {
			curie, ok := ParseCURIE(tc.Input)
			if _a, _e := ok, tc.ExpectedOK; _a != _e {
				t.Fatalf("ok: expected %v, got %v", _e, _a)
			} else if _a, _e := curie, tc.ExpectedCURIE; _a != _e {
				t.Fatalf("curie: expected %v, got %v", _e, _a)
			}
		})
	}
}
