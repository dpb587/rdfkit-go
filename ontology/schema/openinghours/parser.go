package openinghours

import (
	"fmt"
	"strings"
	"time"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
)

type ParseTextResult struct {
	Data PropertyValue

	ErrFatal error
}

// Err will be non-nil if there was a fatal error during parsing.
func (r ParseTextResult) Err() error {
	if r.ErrFatal != nil {
		return r.ErrFatal
	}

	return nil
}

func ParseText(v string) ParseTextResult {
	var res ParseTextResult

	lexicalSplit := strings.Fields(xsdutil.WhiteSpaceCollapse(v))

	var lexicalDayRanges, lexicalTimeRange string

	if len(lexicalSplit) == 2 {
		lexicalDayRanges = lexicalSplit[0]
		lexicalTimeRange = lexicalSplit[1]
	} else if len(lexicalSplit) == 1 {
		if lexicalSplit[0] >= "0" && lexicalSplit[0] <= "9" {
			lexicalTimeRange = lexicalSplit[0]
		} else {
			lexicalDayRanges = lexicalSplit[0]
		}
	}

	for _, lexicalDayRange := range strings.Split(lexicalDayRanges, ",") {
		var ok bool
		var rangeFrom, rangeThru uint8

		if strings.Contains(lexicalDayRange, "-") {
			dayRange := strings.SplitN(lexicalDayRange, "-", 2)
			if len(dayRange) != 2 {
				res.ErrFatal = fmt.Errorf("parse days: unexpected syntax: %s", lexicalDayRange)

				return res
			}

			if rangeFrom, ok = dayTokens[dayRange[0]]; !ok {
				res.ErrFatal = fmt.Errorf("parse days: invalid token: %s", dayRange[0])

				return res
			} else if rangeThru, ok = dayTokens[dayRange[1]]; !ok {
				res.ErrFatal = fmt.Errorf("parse days: invalid token: %s", dayRange[1])

				return res
			}

			if rangeThru < rangeFrom {
				rangeThru += 14
			}
		} else {
			if rangeFrom, ok = dayTokens[lexicalDayRange]; !ok {
				res.ErrFatal = fmt.Errorf("parse days: invalid token: %s", lexicalDayRange)

				return res
			}

			rangeThru = rangeFrom
		}

		for i := rangeFrom; i <= rangeThru; i += 2 {
			switch dayCycle[i : i+2] {
			case "Mo":
				res.Data.DayOfWeekMo = true
			case "Tu":
				res.Data.DayOfWeekTu = true
			case "We":
				res.Data.DayOfWeekWe = true
			case "Th":
				res.Data.DayOfWeekTh = true
			case "Fr":
				res.Data.DayOfWeekFr = true
			case "Sa":
				res.Data.DayOfWeekSa = true
			case "Su":
				res.Data.DayOfWeekSu = true
			}
		}
	}

	if len(lexicalTimeRange) > 0 {
		timesSplit := strings.SplitN(lexicalTimeRange, "-", 2)

		if len(timesSplit) != 2 {
			res.ErrFatal = fmt.Errorf("parse times: invalid syntax: %s", lexicalTimeRange)

			return res
		}

		{
			opens, err := time.Parse("15:04", timesSplit[0])
			if err != nil {
				res.ErrFatal = fmt.Errorf("parse times: parse opens: invalid syntax: %s", timesSplit[0])

				return res
			}

			res.Data.OpensTime = opens.Format("15:04")
		}

		{
			closes, err := time.Parse("15:04", timesSplit[1])
			if err != nil {
				res.ErrFatal = fmt.Errorf("parse times: parse closes: invalid syntax: %s", timesSplit[1])

				return res
			}

			res.Data.ClosesTime = closes.Format("15:04")
		}
	}

	return res
}
