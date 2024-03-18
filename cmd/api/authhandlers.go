package main

import (
	"filmoteka/domain/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Handle auth here
	switch r.Method {
	case http.MethodPost:
		app.login(w, r)
	default:
		app.errorJSON(w, fmt.Errorf("error method not allowed"), http.StatusMethodNotAllowed)
	}
}

func (app *Config) login(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := app.readJSON(r, w, &user)
	log.Println(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	pwd, role, err := user.GetByEmail()
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(user.Password))
	if err != nil {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// create session
	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	app.SessionManager.Put(r.Context(), "role", role)

	app.writeJSON(w, http.StatusOK, jsonResponse{Error: false, Message: "Logged in", Data: nil})
}
