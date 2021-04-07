package errors

import "fmt"

type base struct {
	msg         string
	stack       []uintptr
	stackString string
}

var _ error = (*base)(nil)
var _ StackTracer = (*base)(nil)
var _ errorStacker = (*base)(nil)


// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(msg string) error {
	return &base{
		msg:   msg,
		stack: callers(2),
	}
}

// Errorf returns an error with the message fmt.Sprintf(format, args...).
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &base{
		msg:   msg,
		stack: callers(2),
	}
}

// implements error
func (b *base) Error() string {
	return b.msg
}

// implements StackTracer
func (b *base) StackTrace() []uintptr {
	return b.stack
}

// implements errorStacker
func (b *base) errorStack() string {
	if b.stackString == "" {
		b.stackString = stackString(b.stack)
	}
	return b.Error() + "\n" + b.stackString
}