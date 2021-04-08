package dlog

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"runtime"
	"time"

	"github.com/labstack/gommon/color"

	"code.qschou.com/qschou/go_center_common/config"
	"code.qschou.com/qschou/go_center_common/gls"
)

//在main中修改
func SetTopic(topic string) {
	//考虑重入
	if _dJsonLog != nil {
		_dJsonLog.Close()
	}
	dir := defaultDir
	if config.Mode == config.ModeDev {
		dir = "/tmp/log/"
	}
	file, err := NewFileBackend(dir, topic+".log_json")
	if err != nil {
		panic(err)
	}
	_dJsonLog = NewDJsonLog(file, topic)
	SetLogger(_dJsonLog)
}

var _dJsonLog *dJsonLog

func GetJsonDLog() *dJsonLog {
	if _dJsonLog == nil {
		SetTopic(defaultTopic)
	}
	return _dJsonLog
}

type dJsonLog struct {
	prefix string
	level  Lvl
	output io.Writer
	levels []string
	color  *color.Color
	dw     *dlogWriter
}

func NewDJsonLog(w io.WriteCloser, topic string) *dJsonLog {
	if len(topic) <= 0 {
		topic = defaultTopic
	}
	l := &dJsonLog{
		level:  INFO,
		prefix: topic,
		color:  color.New(),
	}
	l.initLevels()
	l.dw = NewDlogWriter(w)
	l.SetOutput(l.dw)
	l.SetLevel(INFO)
	return l
}

//kv 应该是成对的 数据, 类似: name,张三,age,10,...
func (p *dJsonLog) logJson(v Lvl, kv ...interface{}) (err error) {
	if v < p.level {
		return nil
	}
	om := NewOrderMap()
	_, file, line, _ := runtime.Caller(3)
	file = p.getFilePath(file)
	traceID, pSpanID, spanID := gls.GetTraceInfo()
	//增加traceId,spanid,pspanid
	om.Set("prefix", p.Prefix())
	om.Set("level", p.levels[v])
	om.Set("cur_time", time.Now().Format(time.RFC3339Nano))
	om.Set("cur_unix_time", time.Now().Unix())
	om.Set("file", file)
	om.Set("line", line)
	om.Set("traceid", traceID)
	om.Set("pspanid", pSpanID)
	om.Set("spanid", spanID)
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	for i := 0; i < len(kv); i += 2 {
		om.Set(fmt.Sprintf("%v", kv[i]), kv[i+1])
	}
	str, _ := json.Marshal(om)
	str = append(str, []byte("\n")...)
	_, err = p.Output().Write(str)
	return
}
func (p *dJsonLog) Debug(kv ...interface{}) {
	p.logJson(DEBUG, kv...)
}
func (p *dJsonLog) Info(kv ...interface{}) {
	p.logJson(INFO, kv...)
}
func (p *dJsonLog) Warn(kv ...interface{}) {
	p.logJson(WARN, kv...)
}
func (p *dJsonLog) Error(kv ...interface{}) {
	p.logJson(ERROR, kv...)
}
func (p *dJsonLog) getFilePath(file string) string {
	dir, base := path.Dir(file), path.Base(file)
	return path.Join(path.Base(dir), base)
}
func (p *dJsonLog) Close() error {
	if p.dw != nil {
		p.dw.Close()
		p.dw = nil
	}
	return nil
}
func (p *dJsonLog) DebugLog(b bool) {
	if b && config.Mode != config.ModeOnLine && config.Mode != config.ModePre {
		if _dJsonLog != nil {
			_dJsonLog.SetLevel(DEBUG)
		}
	} else {
		if _dJsonLog != nil {
			_dJsonLog.SetLevel(INFO)
		}
	}
}

type Lvl uint8

const (
	DEBUG Lvl = iota + 1
	INFO
	WARN
	ERROR
	OFF
)

func (l *dJsonLog) initLevels() {
	l.levels = []string{
		"-",
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
	}
}

func (l *dJsonLog) Prefix() string {
	return l.prefix
}

func (l *dJsonLog) SetPrefix(p string) {
	l.prefix = p
}

func (l *dJsonLog) Level() Lvl {
	return l.level
}

func (l *dJsonLog) SetLevel(v Lvl) {
	l.level = v
}

func (l *dJsonLog) Output() io.Writer {
	return l.output
}

func (l *dJsonLog) SetOutput(w io.Writer) {
	l.output = w
}

func (l *dJsonLog) Color() *color.Color {
	return l.color
}
