package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
)

// init "database"
var recipeDB = NewBreadDB()

func main() {
	router := httprouter.New()
	router.GET("/bread/saveRecipe", handle.With(saveRecipe))
	router.GET("/bread/getRecipe", handle.With(getRecipe))

	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
