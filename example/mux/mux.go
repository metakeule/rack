package main

import (
	"fmt"
	. "gopkg.in/metakeule/goh4.v5/tag"
	"gopkg.in/metakeule/rack.v5/router"
	"gopkg.in/metakeule/rack.v5/wrapper"
	"log"
	"net/http"
)

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

var (
	index, admin, inner *router.Router
)

func init() {
	index = router.New("/",
		wrapper.Around{onlyLog("beforeOuter"), onlyLog("afterOuter")},
		wrapper.Before{http.HandlerFunc(LogUrl)},
		wrapper.Guard{http.HandlerFunc(favicon)},
	)
	index.NotFoundHandler = http.HandlerFunc(notFound)
	admin = index.SubRouter("/admin/", wrapper.Around{onlyLog("beforeAdmin"), onlyLog("afterAdmin")})
	inner = admin.SubRouter("/inner/", wrapper.Around{onlyLog("beforeInner"), onlyLog("afterInner")})

	index.Handle("/a", webAndLog("a"))
	index.Handle("/b", webAndLog("b"))

	admin.Handle("/", webAndLog("admin"))
	admin.Handle("/a", webAndLog("adminA"))
	admin.Handle("/b", webAndLog("adminB"))

	inner.Handle("/a", webAndLog("innerA"))
	inner.Handle("/b", webAndLog("innerB"))

	router.Mount(http.DefaultServeMux, inner, admin, index)
}

func main() {
	fmt.Println("listening at localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("can't start server: ", err.Error())
	}
}
