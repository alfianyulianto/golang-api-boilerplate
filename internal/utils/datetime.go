package utils

import (
	"time"
)

func ParseDateTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	t, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		return time.Time{}
	}
	return t
}

func FormatTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
