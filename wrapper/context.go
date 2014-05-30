package wrapper

import (
	"fmt"
	// "fmt"
	"net/http"
	"reflect"
)

// copied from github.com/metakeule/meta Assoc
// assoc associates targetPtrPtr with srcPtr so that
// targetPtrPtr is a pointer to srcPtr and
// targetPtr and srcPtr are pointing to the same address
func assoc(srcPtr, targetPtrPtr interface{}) {
	reflect.ValueOf(targetPtrPtr).Elem().Set(reflect.ValueOf(srcPtr))
}

// copied from github.com/metakeule/meta newPtr
// returns a reference to a new reference to a new empty value based on Type
func newPtr(ty reflect.Type) interface{} {
	val := reflect.New(ty)
	ref := reflect.New(val.Type())
	ref.Elem().Set(val)
	return ref.Interface()
}

// calls function with params, but doesn't return anything
func call(fn reflect.Value, params ...reflect.Value) {
	fn.Call(params)
}

func HandlerMethod2(fn interface{}) http.Handler {
	fnVal := reflect.ValueOf(fn)
	ty := fnVal.Type().In(0)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := reflect.New(ty)
		wVal := reflect.ValueOf(w)
		val.Elem().Set(wVal)
		call(fnVal, val.Elem(), wVal, reflect.ValueOf(r))
	})
}

func HandlerMethod(fn interface{}) http.Handler {
	fnVal := reflect.ValueOf(fn)
	numIn := fnVal.Type().NumIn()
	typs := make([]reflect.Type, numIn-2)
	for i := 0; i < numIn-2; i++ {
		typs[i] = fnVal.Type().In(i).Elem()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wVal := reflect.ValueOf(w)

		params := make([]reflect.Value, numIn)
		for i := 0; i < numIn-2; i++ {
			target := newPtr(typs[i])
			UnWrap(w, target)
			params[i] = reflect.Indirect(reflect.ValueOf(target))
		}

		params[numIn-2] = wVal
		params[numIn-1] = reflect.ValueOf(r)
		fnVal.Call(params)
	})
}

//func MustUnWrap(src interface{}, target interface{}) {
func MustUnWrap(src http.ResponseWriter, target interface{}) {
	err := UnWrap(src, target)
	if err != nil {
		panic(err.Error())
	}
}

// consider a struct that is a http.ResponseWriter via embedding
// now we want to unwrap this struct to get its properties.
// since the struct we are looking for might not be the src but
// instead itself wrapped inside the ResponseWriter property of
// the given src we will do it recursivly untill be get
// the struct we look for or did not find it
func UnWrap(src http.ResponseWriter, target interface{}) error {
	srcVl := reflect.ValueOf(src)
	var srcIsPtr bool

	//if srcVl.Kind() != reflect.Ptr {
	//return fmt.Errorf("src must be pointer")
	//}

	if srcVl.Kind() == reflect.Ptr {
		srcIsPtr = true
		srcVl = reflect.Indirect(srcVl)
	}
	if srcVl.Kind() != reflect.Struct {
		return fmt.Errorf("src must be a struct or a pointer to a struct")
	}
	tgtVl := reflect.ValueOf(target)
	if tgtVl.Kind() != reflect.Ptr ||
		reflect.Indirect(tgtVl).Kind() != reflect.Ptr ||
		reflect.Indirect(reflect.Indirect(tgtVl)).Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a pointer to a struct: %T\n", target)
	}

	if srcVl.Type() == reflect.Indirect(reflect.Indirect(tgtVl)).Type() {
		if srcIsPtr {
			assoc(src, target)
		} else {
			ref := reflect.New(srcVl.Type())
			ref.Elem().Set(srcVl)
			assoc(ref.Interface(), target)
		}
		return nil
	}

	field := srcVl.FieldByName("ResponseWriter")

	if !field.IsValid() {
		return fmt.Errorf("has no field ResponseWriter: %T", src)
	}

	if field.IsNil() {
		return fmt.Errorf("ResponseWriter of %T is nil", src)
	}

	fkind := field.Elem().Kind()

	if fkind == reflect.Ptr {
		fkind = reflect.Indirect(field.Elem()).Kind()
	}

	if fkind != reflect.Struct {
		return fmt.Errorf(
			"ResponseWriter of %T is no struct, but %T",
			src,
			reflect.Indirect(field.Elem()).Type().String())
	}

	rw, ok := field.Interface().(http.ResponseWriter)
	if !ok {
		return fmt.Errorf("ResponseWriter field is no http.ResponseWriter, but %T", field.Interface())
	}
	return UnWrap(rw, target)
}

type context struct {
	//	Type interface{}
	Type reflect.Type
}

func Context(ty http.ResponseWriter) context {
	vl := reflect.ValueOf(ty)
	if vl.Kind() == reflect.Ptr {
		ptrTarget := reflect.Indirect(vl)
		if ptrTarget.Kind() != reflect.Struct {
			panic("context must be a struct or a pointer to a struct")
		}
		return context{ptrTarget.Type()}
		//fmt.Printf("type is %s, kind is %s\n", vl.Type().String(), vl.Kind().String())
	}
	if vl.Kind() != reflect.Struct {
		panic("context must be a struct or a pointer to a struct")
	}
	return context{vl.Type()}
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
