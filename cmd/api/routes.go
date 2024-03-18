package main

import "net/http"

func (app *Config) routes() http.Handler {
	mux := http.NewServeMux()

	//auth
	mux.HandleFunc("/auth", app.Authenticate)

	//actor CRUD ops !done
	mux.HandleFunc("/actor", app.HandleActors)

	//add or delete movie from the list of movies for an actor
	mux.HandleFunc("/movie/actormovie", app.HandleActorMovie)

	//find movie by name/actors name //!done
	mux.HandleFunc("/movie/findbyname", app.HandleMoviesByName)

	//find movie by id and lisr all actors with the list of movies for each actor
	mux.HandleFunc("/movie/actorswithmovies", app.GetActorsAndMoviesForMovie)

	//movie CRUD ops !done
	mux.HandleFunc("/movie", app.HandleMovies)

	return app.SessionManager.LoadAndSave(mux)
}
