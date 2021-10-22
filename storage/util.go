package storage

import "time"

func ISOString(time time.Time) string {
	return time.UTC().Format("2006-01-02T15:04:05.999Z07:00")
}