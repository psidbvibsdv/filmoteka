package main

import (
	"errors"
	"filmoteka/domain/models"
	"log"
	"net/http"
	"strconv"
)

func (app *Config) HandleActors(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && len(r.URL.Query()) > 0:
		var actor *models.Actor
		var err error
		app.readJSON(r, w, &actor)
		actor.ActorID, err = strconv.Atoi(r.URL.Query().Get("actorid"))
		if err != nil {
			log.Println("Error getting actor, unsupported value", err)
			app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
			return
		}
		res, err := actor.GetByID()
		if err != nil {
			log.Println("Error getting actor", err)
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Actor retrieved", Data: res})

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
		res, err := actor.Update()
		if err != nil {
			log.Println("Error updating actor", err)
			app.errorJSON(w, errors.New("error updating actor"), http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Actors updated", Data: res})

		//handler for a delete request
	case r.Method == http.MethodDelete:
		var err error
		r.URL.Query().Get("actorid")
		actor := &models.Actor{}
		actor.ActorID, err = strconv.Atoi(r.URL.Query().Get("actorid"))
		if err != nil {
			log.Println("Error deleting actor, unsupported value", err)
			app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
			return
		}
		err = actor.Delete()
		if err != nil {
			log.Println("Error deleting actor", err)
			app.errorJSON(w, errors.New("error creating actor"), http.StatusInternalServerError)
			return
		}
		app.writeJSON(w, http.StatusCreated, jsonResponse{Error: false, Message: "Actor successfully deleted"})

	}

}
