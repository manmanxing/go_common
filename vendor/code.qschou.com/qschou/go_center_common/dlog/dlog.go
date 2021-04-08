package dlog

import (
	"fmt"
	"io"
	"path"
	"runtime"

	"github.com/labstack/gommon/log"

	"code.qschou.com/qschou/go_center_common/config"
	"code.qschou.com/qschou/go_center_common/gls"
)

const defaultDir = "/home/golanger/log/"

//在main中修改
//func SetTopic(topic string) {
//	//考虑重入
//	if _dLog != nil {
//		_dLog.Close()
//	}
//	dir := defaultDir
//	if config.Mode == config.ModeDev {
//		dir = "/tmp/log/"
//	}
//	file, err := NewFileBackend(dir, topic+".log")
//	if err != nil {
//		panic(err)
//	}
//	_dLog = NewDLog(file, topic)
//	SetLogger(_dLog)
//}

var _dLog *dLog

func GetDLog() *dLog {
	if _dLog == nil {
		SetTopic(defaultTopic)
	}
	return _dLog
}

type dLog struct {
	*log.Logger
	dw *dlogWriter
}

const defaultTopic = "default_topic"

//const defaultHeader = `${prefix} ${level} ${time_rfc3339} ${short_file} ${line}`
const defaultHeader = `${prefix} ${level} ${time_rfc3339}`

func NewDLog(w io.WriteCloser, topic string) *dLog {
	if len(topic) <= 0 {
		topic = defaultTopic
	}
	ret := &dLog{
		Logger: log.New(topic),
	}
	ret.dw = NewDlogWriter(w)
	ret.SetOutput(ret.dw)
	ret.SetHeader(defaultHeader)
	ret.SetLevel(log.INFO)
	switch config.Mode {
	case config.ModeDev, config.ModeQa:
		//ret.SetLevel(log.DEBUG)
		ret.EnableColor()
	case config.ModePre, config.ModeOnLine:
		//ret.SetLevel(log.INFO)
		ret.DisableColor()
	default:
		panic(fmt.Sprintf("unknown,config.Mode=%v", config.Mode))
	}
	return ret
}

//kv 应该是成对的 数据, 类似: name,张三,age,10,...
func (p *dLog) logStr(kv ...interface{}) string {
	_, file, line, _ := runtime.Caller(3)
	file = p.getFilePath(file)
	traceID, pSpanID, spanID := gls.GetTraceInfo()
	//增加traceId,spanid,pspanid
	pre := []interface{}{"traceid", traceID, "pspanid", pSpanID, "spanid", spanID}
	kv = append(pre, kv...)
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	strFmt := "%s %d "
	args := []interface{}{file, line}
	for i := 0; i < len(kv); i += 2 {
		strFmt += "[%v=%+v]"
		args = append(args, kv[i], kv[i+1])
	}
	str := fmt.Sprintf(strFmt, args...)
	return str
}
func (p *dLog) Debug(kv ...interface{}) {
	p.Debugf("", p.logStr(kv...))
}
func (p *dLog) Info(kv ...interface{}) {
	p.Infof("", p.logStr(kv...))
}
func (p *dLog) Warn(kv ...interface{}) {
	p.Warnf("", p.logStr(kv...))
}
func (p *dLog) Error(kv ...interface{}) {
	p.Errorf("", p.logStr(kv...))
}
func (p *dLog) getFilePath(file string) string {
	dir, base := path.Dir(file), path.Base(file)
	return path.Join(path.Base(dir), base)
}
func (p *dLog) Close() error {
	if p.dw != nil {
		p.dw.Close()
		p.dw = nil
	}
	return nil
}
func (p *dLog) DebugLog(b bool) {
	if b && config.Mode != config.ModeOnLine && config.Mode != config.ModePre {
		if _dLog != nil {
			GetDLog().SetLevel(log.DEBUG)
		}

	} else {
		if _dLog != nil {
			GetDLog().SetLevel(log.INFO)
		}
	}
}
