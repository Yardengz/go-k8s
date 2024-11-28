package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type Animal interface {
	Speak() string
}

type Dog struct {
	Name string
}

type Cat struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name + ": woof"
}

func (c Cat) Speak() string {
	return c.Name + ": meow"
}

// In-memory storage for animals
var animals = []Animal{
	Dog{Name: "Rey"},
	Cat{Name: "Mitzi"},
	Dog{Name: "Nala"},
	Cat{Name: "Mutzi"},
	Dog{Name: "Shendy"},
}

func main() {
	r := chi.NewRouter()

	// Middleware for logging and recovering from panics
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check routes for Kubernetes probes
	r.Get("/health", health) // Liveness probe endpoint

	// Routes
	r.Get("/animals", getAnimals)
	r.Post("/animals", addAnimal)
	r.Get("/animals/{name}", getAnimalByName)

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

// Liveness check: This endpoint checks if the application is alive
func health(w http.ResponseWriter, r *http.Request) {
	// You could add more checks here, e.g., if your DB or external service is reachable
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK) // Application is alive
}

// Handlers

// Get all animals
func getAnimals(w http.ResponseWriter, r *http.Request) {
	var response []map[string]string
	for _, animal := range animals {
		switch a := animal.(type) {
		case Dog:
			response = append(response, map[string]string{"type": "dog", "name": a.Name, "sound": a.Speak()})
		case Cat:
			response = append(response, map[string]string{"type": "cat", "name": a.Name, "sound": a.Speak()})
		}
	}

	json.NewEncoder(w).Encode(response)
}

// Add a new animal
func addAnimal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	switch req.Type {
	case "dog":
		animals = append(animals, Dog{Name: req.Name})
	case "cat":
		animals = append(animals, Cat{Name: req.Name})
	default:
		http.Error(w, "Invalid animal type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Animal added"))
}

// Get an animal by name
func getAnimalByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	for _, animal := range animals {
		switch a := animal.(type) {
		case Dog:
			if a.Name == name {
				json.NewEncoder(w).Encode(map[string]string{"type": "dog", "name": a.Name, "sound": a.Speak()})
				return
			}
		case Cat:
			if a.Name == name {
				json.NewEncoder(w).Encode(map[string]string{"type": "cat", "name": a.Name, "sound": a.Speak()})
				return
			}
		}
	}

	http.Error(w, "Animal not found", http.StatusNotFound)
}
