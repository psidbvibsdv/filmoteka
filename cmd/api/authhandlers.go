package main

import "net/http"

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Handle auth here
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}
