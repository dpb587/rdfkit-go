package xsdtype

import "testing"

func TestMapDuration(t *testing.T) {
	for _, tc := range []struct {
		InputString   string
		OutputLiteral Duration
		Error         string
	}{
		{
			InputString: "P1Y2M3DT4H5M6S",
			OutputLiteral: Duration{
				Years:   1,
				Months:  2,
				Days:    3,
				Hours:   4,
				Minutes: 5,
				Seconds: 6,
			},
		},
		{
			InputString: "-P1Y2M3DT4H5M6S",
			OutputLiteral: Duration{
				Years:    1,
				Months:   2,
				Days:     3,
				Hours:    4,
				Minutes:  5,
				Seconds:  6,
				Negative: true,
			},
		},
		{
			InputString: "P1.1Y2.2M3.3DT4.4H5.5M6.6S",
			OutputLiteral: Duration{
				Years:   1.1,
				Months:  2.2,
				Days:    3.3,
				Hours:   4.4,
				Minutes: 5.5,
				Seconds: 6.6,
			},
		},
		{
			InputString: "P1Y",
			OutputLiteral: Duration{
				Years: 1,
			},
		},
		{
			InputString: "P1M",
			OutputLiteral: Duration{
				Months: 1,
			},
		},
		{
			InputString: "P1D",
			OutputLiteral: Duration{
				Days: 1,
			},
		},
		{
			InputString: "PT1H",
			OutputLiteral: Duration{
				Hours: 1,
			},
		},
		{
			InputString: "PT1M",
			OutputLiteral: Duration{
				Minutes: 1,
			},
		},
		{
			InputString: "PT1S",
			OutputLiteral: Duration{
				Seconds: 1,
			},
		},
	} {
		t.Run(tc.InputString, func(t *testing.T) {
			output, err := MapDuration(tc.InputString)
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil {
				if err.Error() != tc.Error {
					t.Errorf("unexpected error: %s", err)
				}
			} else if _e, _a := tc.OutputLiteral, output; _e != _a {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			}
		})
	}
}
