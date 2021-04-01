package errors


func Cause(err error)error  {
	var (
		causer Causer
		ok     bool
	)
	//当找到第一个不是实现这个接口的Error就返回。
	for err != nil {
		causer, ok = err.(Causer)
		if !ok {
			break
		}
		err = causer.Cause()
	}
	return err
}

func String(err error) string {
	if err == nil {
		return ""
	}
	v, ok := err.(StackTracer)
	if !ok {
		return err.Error()
	}
	stack := v.StackTrace()
	if len(stack) == 0 {
		return err.Error()
	}
	if v, ok := err.(errorStacker); ok {
		return v.errorStack()
	}
	return err.Error() + "\n" + stackString(stack)
}