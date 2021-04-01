package gorm

import (
	"database/sql/driver"
	"fmt"
	"github.com/labstack/gommon/log"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

//自定义 gorm log 输出
func NewLogger() logger {
	return logger{}
}

var (
	sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)
)

const (
	timeFormat = "2006-01-02 15:04:05.999999"
)

// Logger default logger
type logger struct {
}

// Print format & print log
func (l logger) Print(values ...interface{}) {
	if len(values) > 1 {
		level := values[0]
		source := fmt.Sprintf("%v", values[1])

		if level == "sql" {
			// duration
			duration := values[2].(time.Duration)
			args := values[4].([]interface{})
			// sql
			var sql string
			var formattedValues = make([]string,0)
			for _, value := range args {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format(timeFormat)))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}

			var formattedValuesLength = len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
			l.sqlLog(sql, args, duration, path.Base(source))

		} else {
			err := values[2]
			log.Error("source", source, "err", err)
		}
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
func (l logger) sqlLog(sql string, args []interface{}, dur time.Duration, source string) {
	argsArr := make([]string, len(args))
	for k, v := range args {
		argsArr[k] = fmt.Sprintf("%v", v)
	}
	argsStr := strings.Join(argsArr, ",")
	//对于超时的,统一打warn日志
	if dur > (time.Millisecond * 400) {
		log.Warn("sql", sql, "args_detal", argsStr, "dur(ms)", int64(dur/time.Millisecond), "source", source)
	} else {
		log.Debug("sql", sql, "args_detal", argsStr, "dur(ms)", int64(dur/time.Millisecond), "source", source)
	}
}