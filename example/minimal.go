package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
)

func Hello(r *http.Request, _ httprouter.Params) handle.Response {
	return handle.Response{
		StatusCode: http.StatusOK,
		Body:       strings.NewReader("Hello, World!"),
	}
}

func Empty(r *http.Request, _ httprouter.Params) handle.Response {
	return handle.Response{}
}

func main() {
	router := httprouter.New()
	router.GET("/", handle.With(Hello))
	router.GET("/bad", handle.With(Empty))

	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
