package models

// Recipe date Model
type Recipe struct {
	ID          string
	Author      string
	Description string
	Images      Images
	Ingredients []Ingredient
}
