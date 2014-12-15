package main

import (
	"fmt"
	"reflect"
)

type x struct {
	Y string
}

func (xx *x) Print(key string) {
	fmt.Printf("%s: %#v\n", key, xx)
}

// a must be a pointer to a value
// and b must be a pointer to a pointer to a value of the same type as the value a points to
func assoc(a, b interface{}) {
	reflect.ValueOf(b).Elem().Set(reflect.ValueOf(a))
}

func newPtr(ty reflect.Type) interface{} {
	val := reflect.New(ty)
	val.Elem().FieldByName("Y").SetString("newPtr")
	return val.Interface()
}

func newPtrPtr(ty reflect.Type) interface{} {
	val := reflect.New(ty)
	val.Elem().FieldByName("Y").SetString("newPtrPtr")
	f := reflect.New(val.Type())
	f.Elem().Set(val)
	return f.Interface()
}

func main() {
	x1 := &x{"hi"}
	x2 := &x{}

	assoc(x1, &x2)

	fmt.Printf("x2: %#v\n", x2)

	x1.Y = "ho"

	fmt.Printf("x2: %#v\n", x2)

	p := newPtr(reflect.TypeOf(x{})).(*x)

	p.Print("p")

	t := newPtrPtr(reflect.TypeOf(x{})).(**x)

	fmt.Printf("t: %#v\n", *t)

	//fmt.Printf("p: %#v\n", p)

}
