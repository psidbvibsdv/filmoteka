package actormoviehandlers

import (
	"errors"
	"filmoteka/internal/domain/models"
	"filmoteka/internal/utils"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
)

type ActorMovieHandler struct {
	useCase        actorMovieUseCase
	sessionManager *scs.SessionManager
}

func New(useCase actorMovieUseCase, manager *scs.SessionManager) *ActorMovieHandler {
	return &ActorMovieHandler{
		useCase:        useCase,
		sessionManager: manager,
	}
}

type actorMovieUseCase interface {
	GetActorsForMovie(movieid int) ([]*models.Actor, *models.Movie, error)
	GetMoviesForActor(actorid int) ([]*models.Movie, *models.Actor, error)
	GetActorsAndMoviesForMovie(movieid int) ([]*models.ActorMovies, *models.Movie, error)
	GetMovieByActorName(firstname string, lastname string) ([]*models.MovieWithActor, error)
	AddActorToMovie(actorid int, movieid int) (*models.Actor, *models.Movie, error)
	DeleteActorFromMovie(actorid int, movieid int) (*models.Actor, *models.Movie, error)
}

func (h *ActorMovieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	role := h.sessionManager.GetString(r.Context(), "role")
	if role != "admin" && r.Method != http.MethodGet {
		utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == http.MethodGet:
		h.getMovie(w, r)
	case r.Method == http.MethodPost:
		h.addActorToMovie(w, r)
	case r.Method == http.MethodDelete:
		h.deleteActorFromMovie(w, r)
	}
}

func (h *ActorMovieHandler) deleteActorFromMovie(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ActorID int `json:"actorid"`
		MovieID int `json:"movieid"`
	}
	req := &request{}
	movie := &models.Movie{}
	actor := &models.Actor{}
	err := utils.ReadJSON(r, w, &req)
	log.Println("Request: ", req)
	if err != nil {
		log.Println("Error reading request", err)
		utils.ErrorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
		return
	}
	actor, movie, err = h.useCase.DeleteActorFromMovie(req.ActorID, req.MovieID)
	if err != nil {
		log.Println("Error deleting actor from movie", err)
		utils.ErrorJSON(w, errors.New("error deleting actor from movie"), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: fmt.Sprintf("Actor  '%s'  succesfully deleted from the movie (%s) ", actor.Name, movie.Title), Data: req})

}

func (h *ActorMovieHandler) getMovie(w http.ResponseWriter, r *http.Request) {

	//allowed request parameters
	expectedParams := map[string]bool{
		"firstname": true,
		"lastname":  true,
		"movieid":   true,
		"actorid":   true,
		"action":    true,
	}
	for param := range r.URL.Query() {
		if !expectedParams[param] {
			log.Println("Invalid request parameter: ", param)
			utils.ErrorJSON(w, errors.New("invalid request parameter. "), http.StatusBadRequest)
			return
		}
	}
	firstnameParam, firstnameOk := r.URL.Query()["firstname"]
	lastnameParam, lastnameOk := r.URL.Query()["lastname"]
	movieid, idOk := r.URL.Query()["movieid"]
	actorid, actorOK := r.URL.Query()["actorid"]
	action, actionOK := r.URL.Query()["action"]

	//action = getmovies + actorid ==> GetMoviesForActor(actorid int)
	//action = getactors + movieid ==> GetActorsForMovie(movieid int)
	//action = getactorandmovie + movieid ==> GetActorsAndMoviesForMovie(movieid int)
	//firstname + lastname ==> GetMovieByActorName(name string, surname string)

	switch {
	case r.Method == http.MethodGet && actionOK && len(action) > 0:
		if action[0] == "getmovies" && actorOK && len(actorid) > 0 {
			var id int
			utils.StringToInt(w, &id, actorid[0])
			movies, actor, err := h.useCase.GetMoviesForActor(id)
			if errors.Is(err, models.ErrNoRecord) {
				utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting movies", err)
				utils.ErrorJSON(w, errors.New("Error getting movies for actor"), http.StatusInternalServerError)
				return
			}
			utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: fmt.Sprint("Movies retrieved for actor ", actor.Name), Data: movies})
		} else if action[0] == "getactors" && idOk && len(movieid) > 0 {
			var id int
			utils.StringToInt(w, &id, movieid[0])
			actors, movie, err := h.useCase.GetActorsForMovie(id)
			if errors.Is(err, models.ErrNoRecord) {
				utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting actors", err)
				utils.ErrorJSON(w, errors.New("Error getting actors for movie"), http.StatusInternalServerError)
				return
			}

			utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: fmt.Sprintf("Actors retrieved for movie `%s`  (%d)", movie.Title, movie.ReleaseDate.Time.Year()), Data: actors})
		} else if action[0] == "getactorandmovie" && idOk && len(movieid) > 0 {
			var id int
			utils.StringToInt(w, &id, movieid[0])
			//if err != nil {
			//	log.Println("Error converting id to int", err)
			//	utils.ErrorJSON(w, errors.New("invalid id parameter"), http.StatusBadRequest)
			//	return
			//}
			res, movie, err := h.useCase.GetActorsAndMoviesForMovie(id)
			if errors.Is(err, models.ErrNoRecord) {
				utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting result", err)
				utils.ErrorJSON(w, err, http.StatusInternalServerError)
				return
			}
			utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: fmt.Sprint("Actors of the movie ", movie.Title, " and their movies: "), Data: res})
		} else {
			utils.ErrorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		}
	case r.Method == http.MethodGet && !actionOK:
		if (firstnameOk && len(firstnameParam) > 0) && (lastnameOk && len(lastnameParam) > 0) {
			// Fetch movies by the name of the actor
			firstname := firstnameParam[0]
			lastname := lastnameParam[0]

			movies, err := h.useCase.GetMovieByActorName(firstname, lastname)
			if errors.Is(err, models.ErrNoRecord) {
				utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting movies", err)
				utils.ErrorJSON(w, err, http.StatusInternalServerError)
				return
			}

			utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: fmt.Sprint("Movie retrieved for actor ", firstnameParam, " ", lastnameParam), Data: movies})

		} else if firstnameOk && len(firstnameParam) > 0 {
			firstname := firstnameParam[0]
			lastname := ""
			//movie := &models.Movie{}
			movies, err := h.useCase.GetMovieByActorName(firstname, lastname)
			if errors.Is(err, models.ErrNoRecord) {
				utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("Error getting movies", err)
				utils.ErrorJSON(w, err, http.StatusInternalServerError)
				return
			}
			utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: fmt.Sprint("Movie retrieved for actor ", firstnameParam, " ", lastnameParam), Data: movies})

		} else if lastnameOk && len(lastnameParam) > 0 {
			firstname := ""
			lastname := lastnameParam[0]
			//movie := &models.Movie{}
			movies, err := h.useCase.GetMovieByActorName(firstname, lastname)
			if errors.Is(err, models.ErrNoRecord) {
				utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
				return

			} else if err != nil {
				log.Println("Error getting movies", err)
				utils.ErrorJSON(w, err, http.StatusInternalServerError)
				return
			}
			utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: fmt.Sprint("Movie retrieved for actor ", firstnameParam, " ", lastnameParam), Data: movies})
		} else {
			utils.ErrorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		}
	}
}

func (h *ActorMovieHandler) addActorToMovie(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ActorID int `json:"actorid"`
		MovieID int `json:"movieid"`
	}
	movie := &models.Movie{}
	actor := &models.Actor{}
	req := &request{}
	err := utils.ReadJSON(r, w, &req)
	log.Println("Movie: ", req)
	if err != nil {
		log.Println("Error reading request", err)
		utils.ErrorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
		return
	}
	actor, movie, err = h.useCase.AddActorToMovie(req.ActorID, req.MovieID)
	if err != nil {
		log.Println("Error adding actor to movie", err)
		utils.ErrorJSON(w, errors.New("error adding actor to movie"), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.JsonResponse{Error: false, Message: fmt.Sprintf("Actor  '%s'  succesfully added to the movie (%s) ", actor.Name, movie.Title), Data: req})

}
