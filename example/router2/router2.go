package main

import (
	"fmt"
	. "gopkg.in/metakeule/goh4.v5/tag"
	// "gopkg.in/metakeule/rack.v5"
	router "gopkg.in/metakeule/rack.v5/router2"
	"gopkg.in/metakeule/rack.v5/wrapper"
	"log"
	"net/http"
)

func init() {
	inner := router.New(wrapper.Around{onlyLog("beforeInner"), onlyLog("afterInner")})
	inner.GET("/a", webAndLog("innerA"))
	inner.GET("/b", webAndLog("innerB"))

	admin := router.New(wrapper.Around{onlyLog("beforeAdmin"), onlyLog("afterAdmin")})
	admin.GET("/", webAndLog("admin"))
	admin.GET("/a", webAndLog("adminA"))
	admin.GET("/b", webAndLog("adminB"))
	admin.GET("/inner", inner)

	other := router.New(wrapper.Around{onlyLog("beforeOther"), onlyLog("afterOther")})
	other.GET("/o", webAndLog("o"))

	index := router.New(
		wrapper.Around{onlyLog("beforeOuter"), onlyLog("afterOuter")},
		wrapper.Before{http.HandlerFunc(LogUrl)},
		wrapper.Guard{http.HandlerFunc(favicon)},
	)
	index.NotFound = http.HandlerFunc(notFound)
	index.GET("/a", webAndLog("a"))
	index.GET("/b", webAndLog("b"))
	index.GET("/admin", admin)
	index.GET("/other", other)

	router.MustMount("/index", index)
}

func main() {
	fmt.Println("listening at localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("can't start server: ", err.Error())
	}
}

type (
	webAndLog string
	onlyLog   string
)

func (ø webAndLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(string(ø))
	P(string(ø)).WriteTo(w)
}

func (ø onlyLog) ServeHTTP(w http.ResponseWriter, r *http.Request) { log.Println(string(ø)) }

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	HTML5(
		HEAD(TITLE("404 - Page not found")),
		BODY(
			H1("404 - Page not found!"), P(r.URL))).WriteTo(w)
}

func LogUrl(w http.ResponseWriter, r *http.Request) { log.Println(r.URL) }

func favicon(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/favicon.ico" {
		w.Write([]byte(""))
	}
}
