package main

import (
	"errors"
	"filmoteka/domain/models"
	"log"
	"net/http"
	"strconv"
)

func (app *Config) HandleActors(w http.ResponseWriter, r *http.Request) {
	role := app.SessionManager.GetString(r.Context(), "role")
	if role == "user" && r.Method != http.MethodGet ||
		role != "admin" {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == http.MethodGet && len(r.URL.Query()) > 0:
		app.getActorsById(w, r)
	case r.Method == http.MethodGet:
		app.getAllActors(w, r)
	case r.Method == http.MethodPost:
		app.createActor(w, r)
	case r.Method == http.MethodPatch:
		app.updateActor(w, r)
	case r.Method == http.MethodDelete:
		app.deleteActor(w, r)
	default:
		app.errorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	}

}

func (app *Config) getActorsById(w http.ResponseWriter, r *http.Request) {
	expectedParams := map[string]bool{
		"id": true,
	}

	for param := range r.URL.Query() {
		if !expectedParams[param] {
			log.Println("Invalid request parameter: ", param)
			app.errorJSON(w, errors.New("invalid request parameter "), http.StatusBadRequest)
			return
		}
	}

	var actor *models.Actor
	var err error
	app.readJSON(r, w, &actor)
	actor.ActorID, err = strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Error getting actor, invalid value", err)
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
}

func (app *Config) createActor(w http.ResponseWriter, r *http.Request) {
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
}

func (app *Config) getAllActors(w http.ResponseWriter, r *http.Request) {
	// Handle get request ==> return all actors
	actors, err := app.Models.Actor.GetAll()
	if errors.Is(err, models.ErrNoRecord) {
		app.errorJSON(w, models.ErrNoRecord, http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error getting actors", err)
		app.errorJSON(w, errors.New("Error getting actors"), http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Actors retrieved", Data: actors})
}

func (app *Config) updateActor(w http.ResponseWriter, r *http.Request) {
	actor := &models.Actor{}
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
}

func (app *Config) deleteActor(w http.ResponseWriter, r *http.Request) {
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
	actor := &models.Actor{}
	actor.ActorID, err = strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Error deleting actor, invalid value", err)
		app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}
	err = actor.Delete()

	if err != nil {
		log.Println("Error deleting actor", err)
		app.errorJSON(w, errors.New("error deleting actor"), http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusCreated, jsonResponse{Error: false, Message: "Actor successfully deleted"})
}
