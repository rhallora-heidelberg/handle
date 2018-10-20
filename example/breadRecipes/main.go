package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
	"log"
	"net/http"
)

type doughRecipe struct {
	Name     string
	GFlour   int
	GSalt    int
	GStarter int
	GWater   int
}

var recipeDB map[string]doughRecipe

func main() {
	// init "database"
	recipeDB = make(map[string]doughRecipe)

	router := httprouter.New()
	router.GET("/bread/saveRecipe", handle.With(saveRecipe))
	router.GET("/bread/getRecipe", handle.With(getRecipe))

	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
