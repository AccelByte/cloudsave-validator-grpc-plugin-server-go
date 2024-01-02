// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

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
