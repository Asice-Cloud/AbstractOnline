package utils

import (
	"fmt"
	"reflect"
)

type ErrorHolder struct {
	err error
}

func Throw(err error) {
	panic(ErrorHolder{err: err})
}

type FinallyHandler interface {
	Finally(handlers ...func())
}

type CatchHandler interface {
	Catch(e error, handler func(err error)) CatchHandler
	CatchAll(handler func(err error)) FinallyHandler
	FinallyHandler
}

type catchHandler struct {
	err      error
	hasCatch bool
}

func Try(f func()) CatchHandler {
	t := &catchHandler{}
	defer func() {
		if r := recover(); r != nil {
			if eh, ok := r.(ErrorHolder); ok {
				t.err = eh.err
			} else if err, ok := r.(error); ok {
				t.err = err
			} else {
				t.err = fmt.Errorf("unknown panic: %v", r)
			}
			t.hasCatch = false
		}
	}()
	f()
	return t
}

func (t *catchHandler) Catch(e error, handler func(err error)) CatchHandler {
	if t.err != nil && !t.hasCatch && reflect.TypeOf(t.err) == reflect.TypeOf(e) {
		handler(t.err)
		t.hasCatch = true
	}
	return t
}

func (t *catchHandler) CatchAll(handler func(err error)) FinallyHandler {
	if t.err != nil && !t.hasCatch {
		handler(t.err)
	}
	return t
}

func (t *catchHandler) Finally(handlers ...func()) {
	for _, handler := range handlers {
		handler()
	}
}

/*
type MyError struct {
	Message string
}

func (e MyError) Error() string {
	return e.Message
}

func riskyOperation() error {
	return MyError{"Something went wrong"}
}

func main() {
	Try(func() {
		xx := riskyOperation()
		throw(xx)
	}).CatchAll(func(err error) {
		fmt.Println("Caught a panic or an unspecified error:", err)
	}).Finally(func() {
		fmt.Println("Operation attempted")
	})
}
*/
