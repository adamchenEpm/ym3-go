package common

import "time"

func TimeToStr(val time.Time) string {
	return val.Format("2006-01-02 15:04:05")
}
