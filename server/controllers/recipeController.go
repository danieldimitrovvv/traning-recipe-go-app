package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type recipeHandlers struct {
	sync.Mutex
	store map[string]Recipe
}

func (h *recipeHandlers) recipes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *recipeHandlers) get(w http.ResponseWriter, r *http.Request) {
	recipes := make([]Recipe, len(h.store))

	h.Lock()
	i := 0
	for _, recipe := range h.store {
		recipes[i] = recipe
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(recipes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *recipeHandlers) getRandomRecipe(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.store))
	h.Lock()
	i := 0
	for id := range h.store {
		ids[i] = id
		i++
	}
	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}

	w.Header().Add("location", fmt.Sprintf("/recipes/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *recipeHandlers) getRecipe(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")

	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		h.getRandomRecipe(w, r)
		return
	}

	h.Lock()
	coaster, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(coaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *recipeHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var recipe Recipe
	err = json.Unmarshal(bodyBytes, &recipe)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	recipe.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[recipe.ID] = recipe
	defer h.Unlock()
}

func newRecipe() *recipeHandlers {
	return &recipeHandlers{
		store: map[string]Recipe{
			"id1": Recipe{
				Author:      "1",
				Description: "Recipe 1",
				Images: Images{
					All:  []string{"https://firebasestorage.googleapis.com/v0/b/training-recipes-app.appspot.com/o/recipes%2Fbanana_bread%2F91531240_814557239036215_6313265375777652736_n.jpg?alt=media&token=8abd1674-0b26-442b-99dd-021424865db1", "https://firebasestorage.googleapis.com/v0/b/training-recipes-app.appspot.com/o/recipes%2Fbanana_bread%2F91593159_158930561987270_4995466201300729856_n.jpg?alt=media&token=1801ca69-3980-4ccf-84c7-603e92103eaf"},
					Main: "https://firebasestorage.googleapis.com/v0/b/training-recipes-app.appspot.com/o/recipes%2Fbanana_bread%2F91365712_750454092406226_6858959026776965120_n.jpg?alt=media&token=23bb4408-12fb-4c02-a2d4-b5924f1f9dad",
				},
				Ingredients: []Ingredient{
					Ingredient{
						Name:     "banani",
						Quantity: "2 br",
					},
					Ingredient{
						Name:     "orizovo brashno",
						Quantity: "190gr",
					},
				},
			},
		},
	}
}
