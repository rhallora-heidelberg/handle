package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
	respondWith "github.com/rhallora-heidelberg/handle/respond-with"
)

// e.g. "http://localhost:8000/bread/getRecipe?name=PompeiiSourdough"
func getRecipe(r *http.Request, _ httprouter.Params) handle.Response {
	// parse input
	name := r.URL.Query().Get("name")

	// validate input
	if name == "" {
		respondWith.Errorf(http.StatusBadRequest, "must specify a recipe name")
	}

	// attempt to retrieve recipe
	recipe, err := recipeDB.Get(name)
	if err != nil {
		return respondWith.Errorf(http.StatusBadRequest, err.Error())
	}

	// return recipe
	return respondWith.JSONOrError(recipe)
}
