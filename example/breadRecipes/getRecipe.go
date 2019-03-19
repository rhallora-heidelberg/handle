package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
	"github.com/rhallora-heidelberg/handle/respondwith"
)

// e.g. "http://localhost:8000/bread/getRecipe?name=PompeiiSourdough"
func getRecipe(r *http.Request, _ httprouter.Params) handle.Response {
	// parse input
	name := r.URL.Query().Get("name")

	// validate input
	if name == "" {
		respondwith.Errorf(http.StatusBadRequest, "must specify a recipe name")
	}

	// attempt to retrieve recipe
	recipe, err := recipeDB.Get(name)
	if err != nil {
		return respondwith.Errorf(http.StatusBadRequest, err.Error())
	}

	// return recipe
	return respondwith.JSONOrError(recipe)
}
