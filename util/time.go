package util

import "time"

const (
	TimeShortFormat_ = "2006-01-02"
	TimeLongFormat_  = "2006-01-02 15:04:05"
	TimeLongFormat   = "20060102150405"
	TimeShortFormat  = "20060102"
)

var BeijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)
