package config

import "code.qschou.com/qschou/go_center_common/beacon/etcd"

const (
	ModeDev    = "dev"
	ModeQa     = "qa"
	ModePre    = "pre"
	ModeOnLine = "online"
)

var Mode = ModeOnLine

func init() {
	value, err := etcd.Get("root/config/common/mode")
	if err != nil {
		panic(err)
	}
	switch value {
	case ModeDev, ModeQa, ModePre, ModeOnLine:
		Mode = value
	default:
		panic("unknown mode: " + value)
	}
}
