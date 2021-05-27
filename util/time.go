package util

import "time"

const (
	TimeFormatDate_       = "2006-01-02"
	TimeFormatDate        = "20060102"
	TimeFormatTime        = "15:04:05"
	TimeFormatDateTime    = "2006-01-02 15:04:05"
	TimeLongFormat        = "20060102150405"
	TimeFormatUTCDateTime = "2006-01-02T15:04:05Z"
)

var BeijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)
