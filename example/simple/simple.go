package main

import (
	"fmt"
	"github.com/gorilla/mux"
	. "gopkg.in/metakeule/goh4.v5/tag"
	. "gopkg.in/metakeule/goh4.v5/tag/short"
	"gopkg.in/metakeule/rack.v5"
	"gopkg.in/metakeule/rack.v5/wrapper"
	"log"
	"net/http"
	"regexp"
	"runtime"
)

var (
	login  = regexp.MustCompile("/login.*")
	err    = regexp.MustCompile("/error.*")
	router = mux.NewRouter()
)

func main() {
	// a set of middlewares - the first a nearer to the app
	middlewareSet := []rack.Wrapper{
		// catch errors and handle them with errorHandler
		wrapper.Catch(errorHandler),

		// guard routes with mustLogin
		wrapper.Guard{http.HandlerFunc(mustLogin)},
	}

	// a new rack for the app
	Rack := rack.New(http.HandlerFunc(app))

	// wrap the app call in a around helper for the layout
	Rack.Wrap(wrapper.Around{Always(`<!DOCTYPE HTML><body style="background-color:yellow;"><div>Header</div>`), Always(`<p>Footer...</p></body>`)})

	// wrap the app by some middlewares defined above
	Rack.Wrap(middlewareSet...)

	/*
	   the middleware order of a request is now:

	   Guard -> Catch -> Around -> app
	*/

	// run the rack
	http.ListenAndServe(":8080", Rack)
}

func init() {
	router.HandleFunc("/", list)
	router.Handle("/somewhere", Always("somewhere with layout"))
}

func app(w http.ResponseWriter, r *http.Request) {
	if err.MatchString(r.URL.Path) {
		panic("this is a panic")
	}
	router.ServeHTTP(w, r)
}

func mustLogin(w http.ResponseWriter, r *http.Request) {
	if login.MatchString(r.URL.Path) {
		w.Write([]byte("Login required!!"))
	}
}

func errorHandler(p interface{}, w http.ResponseWriter, r *http.Request) {
	str := callback(4)
	log.Printf("Error: %v\n%s\n\n", p, str)
	writeError(p, str, w)
}

func callback(skip int) (str string) {
	str = ""
	for i := 0; i < 40; i++ {
		_, file, line, ok := runtime.Caller(skip + i)
		if ok {
			str += fmt.Sprintf("in %s (%v)\n", file, line)
		}
	}
	return
}

func writeError(m interface{}, details string, w http.ResponseWriter) {
	HTML5(
		BODY(
			H1("Error: "+fmt.Sprintf("%v", m)),
			PRE(details))).WriteTo(w)
}

// writes always the same given string
type Always string

func (ø Always) ServeHTTP(w http.ResponseWriter, r *http.Request) { w.Write([]byte(ø)) }

func list(w http.ResponseWriter, r *http.Request) {
	P("try one of these",
		UL(
			LI(AHref("/login/test", "protected area")),
			LI(AHref("/error/test", "error")),
			LI(AHref("/somewhere", "somewhere")),
		),
	).WriteTo(w)
}
