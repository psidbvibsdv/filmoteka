package main

import (
	"errors"
	"filmoteka/domain/models"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Handle auth here
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func (app *Config) HandleActors(w http.ResponseWriter, r *http.Request) {
	switch {
	//handler for a post request
	case r.Method == http.MethodPost:
		// Handle post request ==> add new actor
		var actor *models.Actor

		err := app.readJSON(r, w, &actor)
		log.Println("Actor: ", actor)
		if err != nil {
			log.Println("Error reading request", err)
			app.errorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
			return
		}

		_, err = actor.Create()
		if err != nil {
			log.Println("Error creating actor", err)
			app.errorJSON(w, errors.New("error creating actor"), http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusCreated, jsonResponse{Error: false, Message: "Actor created", Data: actor})

		//handler for a get request
	case r.Method == http.MethodGet:
		// Handle get request ==> return all actors
		actors, err := app.Models.Actor.GetAll()
		if err != nil {
			log.Println("Error getting actors", err)
			app.errorJSON(w, errors.New("error getting actors"), http.StatusInternalServerError)
			return
		}

		app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Actors retrieved", Data: actors})

		//handler for a put request
	case r.Method == http.MethodPatch:
		// Handle put request ==> update actor and return updated entry
		actor := &models.Actor{
			//Gender:      "",
			//DateOfBirth: models.Date{},
			//Name:        "",
		}
		err := app.readJSON(r, w, &actor)
		if err != nil {
			log.Println("Error reading request", err)
			app.errorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
			return
		}
		_, err = actor.Update()
		if err != nil {
			log.Println("Error updating actor", err)
			app.errorJSON(w, errors.New("error updating actor"), http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Actors updated", Data: actor})

	}

}

func (app *Config) HandleMovies(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && len(r.URL.Query()) > 0:
		// Handle getAll request with result sorted according to the method in the query
		//sortby := r.URL.Query().Get("sortby")
		//actor := r.URL.Query().Get("actor")
		//name := r.URL.Query().Get("name")

	case r.Method == http.MethodPost:
		var movie models.Movie
		app.readJSON(r, w, &movie)
		// Handle post request ==> add new movie

		//TODO add more cases as needed
	}
}
