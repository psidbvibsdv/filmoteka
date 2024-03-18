package main

import (
	"errors"
	"filmoteka/domain/models"
	"fmt"
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
			if errors.Is(err, models.ErrNoRecord) {
				app.errorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting movies", err)
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Movies retrieved", Data: movies})

		} else {
			// Fetch all movies
			movies, err := app.Models.Movie.GetAll("")
			if errors.Is(err, models.ErrNoRecord) {
				app.errorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting movies", err)
				app.errorJSON(w, err, http.StatusInternalServerError)
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

// HandleMoviesByName retrieve movie by the name of the movie/name of the actor
func (app *Config) HandleMoviesByName(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		expectedParams := map[string]bool{
			"firstname": true,
			"lastname":  true,
			"moviename": true,
		}

		for param := range r.URL.Query() {
			if !expectedParams[param] {
				log.Println("Invalid request parameter: ", param)
				app.errorJSON(w, errors.New("invalid request parameter. "), http.StatusBadRequest)
				return
			}
		}

		firstnameParam, firstnameOk := r.URL.Query()["firstname"]
		lastnameParam, lastnameOk := r.URL.Query()["lastname"]
		movieParam, movieOk := r.URL.Query()["moviename"]
		movie := &models.Movie{}

		if firstnameOk && movieOk || lastnameOk && movieOk {
			app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)

		} else if firstnameOk && len(firstnameParam) > 0 && lastnameOk && len(lastnameParam) > 0 {
			// Fetch movies by the name of the actor
			firstname := firstnameParam[0]
			lastname := lastnameParam[0]

			movies, err := movie.GetByActorName(firstname, lastname)
			if errors.Is(err, models.ErrNoRecord) {
				app.errorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting movies", err)
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}

			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: fmt.Sprint("Movie retrieved for actor ", firstnameParam, " ", lastnameParam), Data: movies})

		} else if firstnameOk && len(firstnameParam) > 0 {
			firstname := firstnameParam[0]
			lastname := ""
			//movie := &models.Movie{}
			movies, err := movie.GetByActorName(firstname, lastname)
			if errors.Is(err, models.ErrNoRecord) {
				app.errorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting movies", err)
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: fmt.Sprint("Movie retrieved for actor ", firstnameParam, " ", lastnameParam), Data: movies})

		} else if lastnameOk && len(lastnameParam) > 0 {
			firstname := ""
			lastname := lastnameParam[0]
			//movie := &models.Movie{}
			movies, err := movie.GetByActorName(firstname, lastname)
			if errors.Is(err, models.ErrNoRecord) {
				app.errorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return

			} else if err != nil {
				log.Println("Error getting movies", err)
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: fmt.Sprint("Movie retrieved for actor ", firstnameParam, " ", lastnameParam), Data: movies})

		} else if movieOk && len(movieParam) > 0 {
			// Fetch all movies and sort
			movies, err := movie.GetByMovieName(movieParam[0])
			if err != nil {
				log.Println("Error getting movies", err)
				app.errorJSON(w, errors.New("error getting movies"), http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Movies retrieved", Data: movies})

		} else {
			app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		}
	}
}

func (app *Config) GetActorsAndMoviesForMovie(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && len(r.URL.Query()) > 0:
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
		movieid, idOk := r.URL.Query()["id"]
		if idOk {
			movie := &models.Movie{}
			// Fetch movie by id
			id, err := strconv.Atoi(movieid[0])
			if err != nil {
				log.Println("Error converting id to int", err)
				app.errorJSON(w, errors.New("invalid id parameter"), http.StatusBadRequest)
				return
			}
			movie.MovieID = id
			res, err := movie.GetActorsAndMoviesForMovie()
			if err != nil {
				log.Println("Error getting result", err)
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
			movie, err = movie.GetByID()
			if err != nil {
				log.Println("Error getting movie", err)
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: fmt.Sprintf("Actors of the movie '%s' (%d) and their movies: ", movie.Title, movie.ReleaseDate.Time.Year()), Data: res})
		} else {
			app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		}
	}
}

func (app *Config) HandleActorMovie(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		type request struct {
			ActorID int `json:"actorid"`
			MovieID int `json:"movieid"`
		}
		req := &request{}
		err := app.readJSON(r, w, &req)
		log.Println("Movie: ", req)
		if err != nil {
			log.Println("Error reading request", err)
			app.errorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
			return
		}
		err = app.Models.Movie.AddActorToMovie(req.ActorID, req.MovieID)
		if err != nil {
			log.Println("Error adding actor to movie", err)
			app.errorJSON(w, errors.New("error adding actor to movie"), http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusCreated, jsonResponse{Error: false, Message: "Actor added to movie", Data: req})

	case r.Method == http.MethodDelete:
		//TODO: implement delete actor from movie
	}
}
