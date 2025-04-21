package curie

import "testing"

func TestCURIE_String(t *testing.T) {
	for _, tc := range []struct {
		Input              CURIE
		ExpectedString     string
		ExpectedSafeString string
	}{
		{
			Input: CURIE{
				Prefix:    "ex",
				Reference: "resource",
			},
			ExpectedString:     "ex:resource",
			ExpectedSafeString: "[ex:resource]",
		},
		{
			Input: CURIE{
				Safe:      true,
				Prefix:    "ex",
				Reference: "resource",
			},
			ExpectedString:     "[ex:resource]",
			ExpectedSafeString: "[ex:resource]",
		},
		{
			Input: CURIE{
				Prefix:    "",
				Reference: "resource",
			},
			ExpectedString:     ":resource",
			ExpectedSafeString: "[:resource]",
		},
		{
			Input: CURIE{
				Safe:      true,
				Prefix:    "",
				Reference: "resource",
			},
			ExpectedString:     "[:resource]",
			ExpectedSafeString: "[:resource]",
		},
		{
			Input: CURIE{
				DefaultPrefix: true,
				Reference:     "resource",
			},
			ExpectedString:     "resource",
			ExpectedSafeString: "[resource]",
		},
		{
			Input: CURIE{
				Safe:          true,
				DefaultPrefix: true,
				Reference:     "resource",
			},
			ExpectedString:     "[resource]",
			ExpectedSafeString: "[resource]",
		},
	} {
		if _a, _e := tc.Input.String(), tc.ExpectedString; _a != _e {
			t.Fatalf("string: expected %v, got %v", _e, _a)
		} else if _a, _e := tc.Input.SafeString(), tc.ExpectedSafeString; _a != _e {
			t.Fatalf("safe string: expected %v, got %v", _e, _a)
		}
	}
}
