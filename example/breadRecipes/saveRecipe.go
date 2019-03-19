package main

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rhallora-heidelberg/handle"
	respondWith "github.com/rhallora-heidelberg/handle/respond-with"
)

// e.g. "http://localhost:8000/bread/saveRecipe?name=PompeiiSourdough&gFlour=1600&gStarter=200&gWater=950&gSalt=35"
func saveRecipe(r *http.Request, _ httprouter.Params) handle.Response {
	recipe := parseRecipeQuery(r.URL.Query())

	// validate input
	if err := recipe.Validate(); err != nil {
		return respondWith.Errorf(http.StatusBadRequest, err.Error())
	}

	// attempt to store recipe
	if err := recipeDB.Put(recipe); err != nil {
		return respondWith.Errorf(http.StatusBadRequest, err.Error())
	}

	// success
	return respondWith.Stringf("Ok! We saved your recipe under the name '%s'!", recipe.Name)
}

func parseRecipeQuery(qVals url.Values) DoughRecipe {
	recipe := DoughRecipe{
		Name: qVals.Get("name"),
	}

	// ugly, but you get the point
	if n, err := strconv.Atoi(qVals.Get("gFlour")); err == nil {
		recipe.GFlour = n
	}
	if n, err := strconv.Atoi(qVals.Get("gSalt")); err == nil {
		recipe.GSalt = n
	}
	if n, err := strconv.Atoi(qVals.Get("gStarter")); err == nil {
		recipe.GStarter = n
	}
	if n, err := strconv.Atoi(qVals.Get("gWater")); err == nil {
		recipe.GWater = n
	}

	return recipe
}
