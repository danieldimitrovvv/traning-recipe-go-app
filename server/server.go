package main

import (
	"net/http"

func main() {
	recipeHandlers := newRecipe()

	http.HandleFunc("/recipes", recipeHandlers.recipes)
	http.HandleFunc("/recipes/", recipeHandlers.getRecipe)
	err := http.ListenAndServe(":8087", nil)

	if err != nil {
		panic(err)
	}
}
