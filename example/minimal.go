package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
	"log"
	"net/http"
)

func Hello(r *http.Request, _ httprouter.Params) handle.Respond {
	return func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, World!")
	}
}

func Empty(r *http.Request, _ httprouter.Params) handle.Respond {
	return nil
}

func main() {
	router := httprouter.New()
	router.GET("/", handle.With(Hello))
	router.GET("/bad", handle.With(Empty))

	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
