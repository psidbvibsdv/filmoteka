package main

import "net/http"

func (app *Config) routes() http.Handler {
	mux := http.NewServeMux()

	//auth
	mux.HandleFunc("/auth", app.Authenticate)

	//actor CRUD ops !done
	mux.HandleFunc("/actor", app.HandleActors)

	//add or delete movie from the list of movies for an actor
	mux.HandleFunc("/actor/actormovie", app.HandleActorMovies)

	//find movie by name/actors name //!done
	mux.HandleFunc("/movie/findbyname", app.HandleMoviesByName)

	//movie CRUD ops !done
	mux.HandleFunc("/movie", app.HandleMovies)

	return mux
}
