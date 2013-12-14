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

// consider a struct that is a http.ResponseWriter via embedding
// now we want to
func UnWrap(src interface{}, target interface{}) error {
	srcVl := reflect.ValueOf(src)

	if srcVl.Kind() != reflect.Ptr {
		panic("src must be pointer")
	}

	if srcVl.Kind() == reflect.Ptr {
		srcVl = reflect.Indirect(srcVl)
	}
	if srcVl.Kind() != reflect.Struct {
		panic("src must be a struct or a pointer to a struct")
	}
	tgtVl := reflect.ValueOf(target)
	if tgtVl.Kind() != reflect.Ptr {
		fmt.Printf("1. target must be a pointer to a pointer to a struct: %T, kind %s\n", target, tgtVl.Kind())
		panic("1. target must be a pointer to a pointer to a struct")
	}

	if reflect.Indirect(tgtVl).Kind() != reflect.Ptr {
		fmt.Printf("2. target must be a pointer to a pointer to a struct: %T, kind %s\n", target, reflect.Indirect(tgtVl).Kind())
		panic("2. target must be a pointer to a pointer to a struct")
	}

	if reflect.Indirect(reflect.Indirect(tgtVl)).Kind() != reflect.Struct {
		fmt.Printf("3. target must be a pointer to a pointer to a struct: %T, kind %s\n", target, reflect.Indirect(tgtVl).Kind())
		panic("3. target must be a pointer to a pointer to a struct")
	}

	//fmt.Printf("%T vs %T\n", src, target)
	// identical type
	//if srcVl.Type() == reflect.Indirect(tgtVl).Type() {
	if reflect.Indirect(reflect.ValueOf(src)).Type() == reflect.Indirect(reflect.Indirect(tgtVl)).Type() {
		assoc(src, target)

		//*target = *src
		//tgtVl.Elem().Set(reflect.ValueOf(src))

		//if srcVl.CanAddr() && tgtVl.CanSet() {
		//		if reflect.ValueOf(src).Kind() == reflect.Ptr {
		//tgtVl.Set(srcVl.Addr())
		//addr := srcVl.Addr()
		//			reflect.ValueOf(target).Elem().Set(reflect.ValueOf(src).Elem())
		//		} else {
		// fmt.Printf("can't set %T with address %T\n", target, src)
		//			tgtVl.Elem().Set(srcVl)
		//		}

		//tgtVl.Elem().Set(srcVl)
		//if srcVl.CanAddr() {
		//tgtVl.Set(srcVl.Addr())
		//tgtVl.Elem().Set(srcVl)
		/*
			} else {
				return fmt.Errorf(
					"can't address: %T",
					src)
			}
		*/
		return nil
	}

	field := srcVl.FieldByName("ResponseWriter")

	if !field.IsValid() {
		return fmt.Errorf(
			"has no field ResponseWriter: %T",
			src)
	}

	if field.IsNil() {
		return fmt.Errorf(
			"ResponseWriter of %T is nil",
			src)
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

	return UnWrap(field.Interface().(http.ResponseWriter), target)
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
