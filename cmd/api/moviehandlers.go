package main

import (
	"errors"
	"filmoteka/domain/models"
	"log"
	"net/http"
	"strconv"
)

func (app *Config) HandleMovies(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		expectedParams := map[string]bool{
			"id":   true,
			"sort": true,
		}

		for param := range r.URL.Query() {
			if !expectedParams[param] {
				log.Println("Invalid request parameter: ", param)
				app.errorJSON(w, errors.New("invalid request parameter. "), http.StatusBadRequest)
				return
			}
		}
		idParam, idOk := r.URL.Query()["id"]
		sortParam, sortOk := r.URL.Query()["sort"]

		if sortOk && idOk {
			app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		} else if idOk {
			movie := &models.Movie{}
			// Fetch movie by id
			id, err := strconv.Atoi(idParam[0])
			if err != nil {
				log.Println("Error converting id to int", err)
				app.errorJSON(w, errors.New("invalid id parameter"), http.StatusBadRequest)
				return
			}
			movie.MovieID = id
			movie, err = movie.GetByID()
			if err != nil {
				log.Println("Error getting movie", err)
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Movie retrieved", Data: movie})
		} else if sortOk {
			// Fetch all movies and sort
			movies, err := app.Models.Movie.GetAll(sortParam[0])
			if err != nil {
				log.Println("Error getting movies", err)
				app.errorJSON(w, errors.New("error getting movies"), http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Movies retrieved", Data: movies})
		} else {
			// Fetch all movies
			movies, err := app.Models.Movie.GetAll("")
			if err != nil {
				log.Println("Error getting movies", err)
				app.errorJSON(w, errors.New("error getting movies"), http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Movies retrieved", Data: movies})
		}

		//handler for a post request
	case r.Method == http.MethodPost:
		// Handle post request ==> add new actor
		var movie *models.Movie

		err := app.readJSON(r, w, &movie)
		log.Println("Movie: ", movie)
		if err != nil {
			log.Println("Error reading request", err)
			app.errorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
			return
		}

		_, err = movie.Create()
		if err != nil {
			log.Println("Error creating movie", err)
			app.errorJSON(w, errors.New("error creating movie"), http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusCreated, jsonResponse{Error: false, Message: "Actor created", Data: movie})

		//handler for a patch request
	case r.Method == http.MethodPatch:
		// Handle patch request ==> update actor and return updated entry
		movie := &models.Movie{}
		err := app.readJSON(r, w, &movie)
		if err != nil {
			log.Println("Error reading request", err)
			app.errorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
			return
		}
		res, err := movie.Update()
		if err != nil {
			log.Println("Error updating movie", err)
			app.errorJSON(w, errors.New("error updating movie"), http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Actors updated", Data: res})

		//handler for a delete request
	case r.Method == http.MethodDelete:
		expectedParams := map[string]bool{
			"id": true,
		}

		for param := range r.URL.Query() {
			if !expectedParams[param] {
				log.Println("Invalid request parameter: ", param)
				app.errorJSON(w, errors.New("invalid request parameter. "), http.StatusBadRequest)
				return
			}
		}
		var err error
		movie := &models.Movie{}
		movie.MovieID, err = strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Println("Error deleting movie, invalid value", err)
			app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
			return
		}
		err = movie.Delete()
		if err != nil {
			log.Println("Error deleting movie", err)
			app.errorJSON(w, errors.New("error deleting movie"), http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusCreated, jsonResponse{Error: false, Message: "Movie successfully deleted"})

	}
}
