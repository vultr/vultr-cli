package utils

import (
	"strconv"
	"strings"
	"time"
)

const (
	DefaultRuneLengthShort int = 22
	DefaultRuneLengthLong  int = 42

	DefaultDateFormatShort string = time.DateOnly
)

// TruncateString returns a shortened, ellipses dots formatted string for
// compact display in the CLI. Rune length is the number of runes that will be
// preserved before the ellipses dots are appended to the end of the string.
func TruncateString(input string, runeLength int) string {
	if len(input) <= 0 {
		return ""
	}

	strRunes := []rune(input)

	if len(strRunes) <= runeLength {
		return input
	}

	// truncate at length
	strRunes = strRunes[:runeLength]

	var strBuilder strings.Builder
	strBuilder.WriteString(string(strRunes))
	strBuilder.WriteString("...")

	return strBuilder.String()
}

// FormatDateTimeString takes in a date time formatted string and returns the
// short date-only format for compact display on the CLI. It will attempt to
// parse the RFC3339 format first, then go's DateTime format. If neither
// succeed, it returns the input, un-altered.
func FormatDateTimeString(input string) string {
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		t, err = time.Parse(time.DateTime, input)
		if err != nil {
			return input
		}
	}

	return t.Format(DefaultDateFormatShort)
}

// FormatTimestampString takes in a time stamp and returns the date-only
// formatted string.
func FormatTimestampString(ts string) string {
	var tsTime time.Time
	tsInt, err := strconv.Atoi(ts)
	if err != nil {
		tsTime = time.Time{}
	} else {
		tsTime = time.Unix(int64(tsInt), 0)
	}
	return tsTime.Format(DefaultDateFormatShort)
}
