package providers

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	// timezoneMap is used to convert time zone to offset
	timezoneMap = map[string]string{
		"GMT": "+0000", "EST": "-0500", "EDT": "-0400",
		"CST": "-0600", "CDT": "-0500", "MST": "-0700", "MDT": "-0600",
		"PST": "-0800", "PDT": "-0700", "AKDT": "-0800", "CHST": "+1000",
		"HST": "-1000", "AST": "-0400", "SST": "-1100", "AKST": "-0800",
	}
)

// TODO: is there something in stdlib that does this...?

func GetUTCTimeForLayout(timeStr string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return t, errors.Wrapf(err, "Error parsing time string:%v", timeStr)
	}

	// replace the timezone with the offset, otherwise go won't convert to UTC correctly
	timezone := t.Location().String()
	if offset, ok := timezoneMap[timezone]; ok {
		layout = strings.Replace(layout, "MST", "-0700", 1)
		timeStr = strings.Replace(timeStr, timezone, offset, 1)
		t, err = time.Parse(layout, timeStr)
		if err != nil {
			return t, errors.Wrapf(err, "Error parsing time string:%v", timeStr)
		}
	}
	return t.UTC().Truncate(time.Second), nil
}
