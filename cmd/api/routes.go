package main

import "net/http"

func (app *Config) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth", app.Authenticate)

	mux.HandleFunc("/actor", app.HandleActors)

	mux.HandleFunc("/movies", app.HandleMovies)

	// Add more endpoints as needed

	return mux
}
