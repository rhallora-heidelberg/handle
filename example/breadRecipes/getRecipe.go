package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
	"github.com/rhallora-heidelberg/handle/respond-with"
	"net/http"
)

// e.g. "http://localhost:8000/bread/getRecipe?name=PompeiiSourdough"
func getRecipe(r *http.Request, _ httprouter.Params) handle.Respond {
	// parse input
	name := r.URL.Query().Get("name")

	// validate input
	if name == "" {
		respondWith.Errorf(http.StatusBadRequest, "error: must specify a recipe name")
	}

	// attempt to retrieve recipe
	recipe, ok := recipeDB[name]
	if !ok {
		return respondWith.Errorf(http.StatusBadRequest, "error: recipe does not exist")
	}

	// return recipe
	errResponse := respondWith.Errorf(http.StatusInternalServerError, "error: failed to unmarshal recipe")
	return respondWith.JSONOrError(recipe, errResponse)
}
