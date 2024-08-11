package utils

import "time"

func ParseLocalTime(curr time.Time) (time.Time, error) {
	const gmtPlus7 = 7 * 60 * 60
	gmtPlus7Location := time.FixedZone("GMT+7", gmtPlus7)

	return time.Date(
		curr.Year(), curr.Month(), curr.Day(),
		curr.Hour(), curr.Minute(), curr.Second(),
		curr.Nanosecond(), gmtPlus7Location), nil
}
