package server

import "time"

func isSameDate(t1, t2 time.Time) bool {
	t1, t2 = t1.UTC(), t2.UTC()

	if t1.Day() != t2.Day() {
		return false
	}

	if t1.Month() != t2.Month() {
		return false
	}

	if t1.Year() != t2.Year() {
		return false
	}

	return true
}
