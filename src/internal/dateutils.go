package internal

import (
	"log/slog"
	"time"
)

func ParseInTime(dateStr string) (time.Time, error) {
	out, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		slog.Error("Error parsing date string")
		slog.Error(err.Error())
		return time.Time{}, err
	}
	return out, nil
}

func FormatTimeForNewsApi(timeObj time.Time) string {
	return timeObj.Format("2006-01-02T15:04:05")
}
