package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/rhallora-heidelberg/handle/respondwith"

	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
)

// "Hello World" using handle.Response directly
func HelloDirect(r *http.Request, _ httprouter.Params) handle.Response {
	return handle.Response{
		StatusCode: http.StatusOK,
		Body:       strings.NewReader("Hello, World!"),
	}
}

// "Hello World" using respondwith
func HelloSugary(r *http.Request, _ httprouter.Params) handle.Response {
	return respondwith.Stringf("Hello, World!")
}

// basic 404
func NotFound(r *http.Request, _ httprouter.Params) handle.Response {
	return respondwith.StatusCode(http.StatusNotFound)
}

func main() {
	router := httprouter.New()
	router.GET("/hello1", handle.With(HelloDirect))
	router.GET("/hello2", handle.With(HelloSugary))
	router.GET("/hl3releasedate", handle.With(NotFound))

	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
