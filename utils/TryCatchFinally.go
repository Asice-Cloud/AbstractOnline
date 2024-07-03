package utils

import (
	"reflect"
)

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
		defer func() {
			r := recover()
			if r != nil {
				t.err = r.(error)
			}
		}()
		f()
	}()
	return t
}

func (t *catchHandler) RequireCatch() bool {
	if t.hasCatch {
		return false
	}
	if t.err == nil {
		return false
	}
	return true
}

func (t *catchHandler) Catch(e error, handler func(err error)) CatchHandler {
	if !t.RequireCatch() {
		return t
	}
	if reflect.TypeOf(e) == reflect.TypeOf(t.err) {
		handler(t.err)
		t.hasCatch = true
	}
	return t
}

func (t *catchHandler) CatchAll(handler func(err error)) FinallyHandler {
	if !t.RequireCatch() {
		return t
	}
	handler(t.err)
	t.hasCatch = true
	return t
}

func (t *catchHandler) Finally(handlers ...func()) {
	for _, handler := range handlers {
		defer handler()
	}
	err := t.err
	if err != nil && !t.hasCatch {
		panic(err)
	}
}

type Err1 struct {
	error
}
type Err2 struct {
	error
}

/*func main() {
	Try(func() {
		fmt.Println("Try 1 error")
		panic(Err1{error: errors.New("error1")})
	}).Catch(Err1{}, func(err error) {
		println("catch err1", err.Error())
	}).Catch(Err2{}, func(err error) {
		println("catch err2", err.Error())
	}).CatchAll(func(err error) {
		println("catch all")
	}).Finally(func() {
		println("finally 1 done")
	})
}
*/
