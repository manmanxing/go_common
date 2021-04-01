package errors

type StackTracer interface {
	error
	StackTrace() []uintptr
}

type errorStacker interface {
	errorStack() string
}

type Causer interface {
	error
	Cause() error
}