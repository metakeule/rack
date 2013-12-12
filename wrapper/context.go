package wrapper

import (
	// "fmt"
	"net/http"
	"reflect"
)

// calls function with params, but doesn't return anything
func call(fn reflect.Value, params ...reflect.Value) {
	fn.Call(params)
}

func HandlerMethod(fn interface{}) http.Handler {
	fnVal := reflect.ValueOf(fn)
	ty := fnVal.Type().In(0)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := reflect.New(ty)
		wVal := reflect.ValueOf(w)
		val.Elem().Set(wVal)
		call(fnVal, val.Elem(), wVal, reflect.ValueOf(r))
	})
}

type context struct {
	//	Type interface{}
	Type reflect.Type
}

func Context(ty interface{}) context {
	vl := reflect.ValueOf(ty)
	if vl.Kind() == reflect.Ptr {
		return context{reflect.Indirect(reflect.ValueOf(ty)).Type()}
		//fmt.Printf("type is %s, kind is %s\n", vl.Type().String(), vl.Kind().String())
	}
	return context{reflect.ValueOf(ty).Type()}
}

func (c context) Wrap(in http.Handler) (out http.Handler) {
	//ty := reflect.Indirect(reflect.ValueOf(c.Type)).Type()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//val := reflect.New(ty)
		val := reflect.New(c.Type)
		val.Elem().FieldByName("ResponseWriter").Set(reflect.ValueOf(w))
		in.ServeHTTP(val.Interface().(http.ResponseWriter), r)
	})
}
