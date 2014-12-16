package main

import (
	"fmt"
	"gopkg.in/metakeule/rack.v5"
	"gopkg.in/metakeule/rack.v5/wrapper"
	"net/http"
	"os"
)

func app(w http.ResponseWriter, r *http.Request) {
	c := w.(*Context)
	fmt.Fprintf(w, "Firstname: %s, Lastname: %s", c.FirstName, c.LastName)
}

// important, that the struct has anonymous field http.ResponseWriter
type Context struct {
	FirstName, LastName string
	http.ResponseWriter
}

func (c Context) Name() string {
	return fmt.Sprintf("%s %s", c.FirstName, c.LastName)
}

type Namer interface {
	Name() string
}

func set2(c *Context, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	c.LastName = q.Get("lastname")
}

func log2(n Namer, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stdout, "### Name: %s ###\n", n.Name())
}

func (c *Context) set1(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	c.FirstName = q.Get("firstname")
}

func (c *Context) log1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stdout, "### Firstname: %s, Lastname: %s ###\n", c.FirstName, c.LastName)
}

func main() {
	r := rack.New(http.HandlerFunc(app))
	r.Wrap(wrapper.CallBefore{log2})
	r.Wrap(wrapper.CallBefore{(*Context).log1})
	r.Wrap(wrapper.CallBefore{set2})
	r.Wrap(wrapper.CallBefore{(*Context).set1})
	r.Wrap(wrapper.Context{&Context{}})

	http.ListenAndServe(":3131", r)
}
