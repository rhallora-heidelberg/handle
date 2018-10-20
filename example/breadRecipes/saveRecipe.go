package main

import (
	"github.com/rhallora-heidelberg/handle"
	"github.com/rhallora-heidelberg/handle/respond-with"
	"net/http"
	"net/url"
	"strconv"
)

// e.g. "http://localhost:8000/bread/saveRecipe?name=PompeiiSourdough&gFlour=1600&gStarter=200&gWater=950&gSalt=35"
func saveRecipe(r *http.Request) handle.Respond {
	recipe := parseRecipeQuery(r.URL.Query())

	// validate input
	if errResponse := validateRecipe(recipe); errResponse != nil {
		return errResponse
	}

	// attempt to store recipe
	if _, exists := recipeDB[recipe.Name]; exists {
		return respondWith.Errorf(http.StatusBadRequest, "error: recipe already exists, please choose a different name")
	}

	recipeDB[recipe.Name] = recipe

	// success
	return respondWith.Stringf("Ok! We saved your recipe under the name '%s'!", recipe.Name)
}

func parseRecipeQuery(qVals url.Values) doughRecipe {
	recipe := doughRecipe{
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

func validateRecipe(recipe doughRecipe) handle.Respond {
	// blank / zero-value errors
	if recipe.Name == "" {
		return respondWith.Errorf(http.StatusBadRequest, "error: recipe must have a name")
	}

	for _, field := range []int{recipe.GFlour, recipe.GWater, recipe.GSalt, recipe.GStarter} {
		if field < 1 {
			return respondWith.Errorf(http.StatusBadRequest, "error: recipe must have flour, water, salt and starter")
		}
	}

	// special cases
	if recipe.GFlour < recipe.GStarter {
		return respondWith.Errorf(http.StatusBadRequest, "error: less flour than starter; acidity will be very high")
	}

	if recipe.GSalt > recipe.GFlour {
		return respondWith.Errorf(http.StatusBadRequest, "error: more salt than flour; please don't")
	}

	return nil
}
