package main

import (
	"errors"
	"sync"
)

type DoughRecipe struct {
	Name     string
	GFlour   int
	GSalt    int
	GStarter int
	GWater   int
}

func (recipe DoughRecipe) Validate() error {
	// blank / zero-value errors
	if recipe.Name == "" {
		return errors.New("recipe must have a name")
	}

	for _, field := range []int{recipe.GFlour, recipe.GWater, recipe.GSalt, recipe.GStarter} {
		if field < 1 {
			return errors.New("recipe must have flour, water, salt and starter")
		}
	}

	// special cases
	if recipe.GFlour < recipe.GStarter {
		return errors.New("less flour than starter; acidity will be very high")
	}

	if recipe.GSalt > recipe.GFlour {
		return errors.New("more salt than flour; please don't")
	}

	return nil
}

// mock database
type BreadDB struct {
	data map[string]DoughRecipe
	mu   sync.RWMutex
}

func NewBreadDB() *BreadDB {
	return &BreadDB{
		data: make(map[string]DoughRecipe),
		mu:   sync.RWMutex{},
	}
}

func (db *BreadDB) Get(name string) (DoughRecipe, error) {
	db.mu.RLock()
	recipe, ok := db.data[name]
	db.mu.RUnlock()

	if !ok {
		return recipe, errors.New("recipe does not exist")
	}

	return recipe, nil
}

func (db *BreadDB) Put(recipe DoughRecipe) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[recipe.Name]; exists {
		return errors.New("recipe with that name already exists")
	}

	db.data[recipe.Name] = recipe
	return nil
}
