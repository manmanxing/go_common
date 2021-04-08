package dlog

//对外提供统一接口，可自定义替换
//默认使用dlog
type Logger interface {
	Debug(kv ...interface{})
	Info(kv ...interface{})
	Warn(kv ...interface{})
	Error(kv ...interface{})
	Close() error
	DebugLog(b bool)
}

var _dLogger Logger

func SetLogger(l Logger) {
	_dLogger = l
}
func GetLogger() Logger {
	if _dLogger == nil {
		SetLogger(GetJsonDLog())
	}
	return _dLogger
}

func Debug(kv ...interface{}) {
	GetLogger().Debug(kv...)
}
func Info(kv ...interface{}) {
	GetLogger().Info(kv...)
}
func Warn(kv ...interface{}) {
	GetLogger().Warn(kv...)
}
func Error(kv ...interface{}) {
	GetLogger().Error(kv...)
}

//这个方法以后不要用了，请使用Close()
func Flush() error {
	return Close()
}
func Close() error {
	return GetLogger().Close()
}
func DebugLog(b bool) {
	GetLogger().DebugLog(b)
}
