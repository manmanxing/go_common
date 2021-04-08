package alert

import (
	"code.qschou.com/golang/errors"
)

type Error struct {
	basicError       error         // 原始error，alert_msg
	inputError       error         // 传进来的error
	desc             string        // 错误概要描述
	threshold        int64         // 告警阈值，1分钟出现几次
	level            Level         // 告警严重程度，关注度
	kvParams         []interface{} // 自定义参数
	isInputErrorWrap bool          // inputError是否被包过
}

func (i *Error) BasicError() error {
	return i.basicError
}

func (i *Error) SetBasicError(basicError error) {
	i.basicError = basicError
}

func (i *Error) InputError() error {
	return i.inputError
}

func (i *Error) SetInputError(inputError error) {
	i.inputError = inputError
}

func (i *Error) Desc() string {
	return i.desc
}

func (i *Error) SetDesc(desc string) {
	i.desc = desc
}

func (i *Error) Threshold() int64 {
	return i.threshold
}

func (i *Error) SetThreshold(threshold int64) {
	i.threshold = threshold
}

func (i *Error) Level() Level {
	return i.level
}

func (i *Error) SetLevel(level Level) {
	i.level = level
}

func (i *Error) KvParams() []interface{} {
	return i.kvParams
}

func (i *Error) SetKvParams(kvParams []interface{}) {
	i.kvParams = kvParams
}

func (i *Error) IsInputErrorWrap() bool {
	return i.isInputErrorWrap
}

func (i *Error) SetIsInputErrorWrap(isInputErrorWrap bool) {
	i.isInputErrorWrap = isInputErrorWrap
}

// new一个error
//  basicError：原始error
//  desc：错误概要描述
//  threshold：告警阈值，1分钟出现几次
//  level：告警严重程度，关注度
func New(inputError error, desc string, threshold int64, level Level, kvParams ...interface{}) error {
	if inputError == nil {
		return nil
	}
	if len(desc) <= 0 {
		return inputError
	}
	if threshold <= 0 {
		return inputError
	}
	if int64(level) <= 0 || int64(level) > 5 {
		return inputError
	}
	var (
		basicError       error
		isInputErrorWrap bool
	)
	ec, assertCauserOK := inputError.(errors.Causer)
	if assertCauserOK {
		basicError = ec.Cause()
		isInputErrorWrap = true
	} else {
		basicError = inputError
	}
	ae, assertAeOK := basicError.(*Error)
	if assertAeOK {
		if !ae.isInputErrorWrap && assertCauserOK {
			ae.inputError = inputError
			ae.isInputErrorWrap = true
		}
		return ae
	}
	if len(kvParams)%2 != 0 {
		kvParams = append(kvParams, "unknown")
	}
	return &Error{
		basicError:       basicError,
		inputError:       inputError,
		desc:             desc,
		threshold:        threshold,
		level:            level,
		kvParams:         kvParams,
		isInputErrorWrap: isInputErrorWrap,
	}
}

func (i *Error) Error() string {
	return i.BasicError().Error()
}
