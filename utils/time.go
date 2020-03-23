package utils

import "time"

func IsTodayDate(t time.Time) bool {
	now := time.Now()
	return t.Month() == now.Month() && t.Day() == now.Day() && t.Year() == now.Year()
}

func IsYestardayDate(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Month() == yesterday.Month() && t.Day() == yesterday.Day() && t.Year() == yesterday.Year()
}
