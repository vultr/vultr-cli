package utils

import (
	"fmt"
	"strconv"
	"time"
)

func ParseAPITimestamp(ts string) string {
	var tsTime time.Time
	tsInt, err := strconv.Atoi(ts)
	if err != nil {
		tsTime = time.Time{}
	} else {
		tsTime = time.Unix(int64(tsInt), 0)
	}

	tsYear, tsMonth, tsDay := tsTime.Date()
	return fmt.Sprintf("%d-%02d-%02d", tsYear, int(tsMonth), tsDay)
}
