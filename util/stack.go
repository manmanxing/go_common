package util

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
)

/*

使用 stackString 函数打印的结果
logic.Do2
        awesomeProject/app/logic/test2.go:4
api.Do1
        awesomeProject/app/api/test1.go:6
main.main
        awesomeProject/app/main.go:41

使用 stackInfo 函数打印的结果
/Users/xxx/go/src/awesomeProject/app/main.go(46): main.stackInfo
/Users/xxx/go/src/awesomeProject/app/main.go(36): main.main.func1
/usr/local/go/src/runtime/panic.go(965): runtime.gopanic
/usr/local/go/src/runtime/panic.go(191): runtime.panicdivide
/Users/xxx/go/src/awesomeProject/app/logic/test2.go(4): awesomeProject/app/logic.Do2
/Users/xxx/go/src/awesomeProject/app/api/test1.go(6): awesomeProject/app/api.Do1
/Users/xxx/go/src/awesomeProject/app/main.go(40): main.main
/usr/local/go/src/runtime/proc.go(225): runtime.main
/usr/local/go/src/runtime/asm_amd64.s(1371): runtime.goexit
*/

//打印所有的栈信息
func StackInfoAll() string {
	var callers string
	for i := 0; true; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		callers += fmt.Sprintf("%s(%d): %s \n", file, line, fn.Name())
	}
	return callers
}

//这里的 skip 一般 <= 3
func Callers(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	return pcs[:n]
}

//打印指定 skip 层的栈信息
func StackInfo(stack []uintptr) string {
	if len(stack) == 0 {
		return ""
	}
	frames := runtime.CallersFrames(stack)

	var (
		frame    runtime.Frame
		more     bool
		funcName string
		fileName string
		buf      bytes.Buffer
	)
	for {
		frame, more = frames.Next()
		if frame.Function == "runtime.main" {
			break
		}
		if frame.Function == "runtime.goexit" {
			break
		}
		if frame.Function == "" {
			funcName = "unknown_function"
		} else {
			funcName = trimFuncName(frame.Function)
		}
		if frame.File == "" {
			fileName = "unknown_file"
		} else {
			fileName = trimFileName(frame.File)
		}
		if buf.Len() > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(funcName)
		buf.WriteString("\n\t")
		buf.WriteString(fileName)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(frame.Line))

		if !more {
			break
		}
	}
	return buf.String()
}

//去除 src 与 vendor 目录
func trimFileName(name string) string {
	i := strings.Index(name, "/src/")
	if i < 0 {
		return name
	}
	name = name[i+len("/src/"):]
	i = strings.Index(name, "/vendor/")
	if i < 0 {
		return name
	}
	return name[i+len("/vendor/"):]
}

func trimFuncName(name string) string {
	return path.Base(name)
}
