package actorhandlers

import (
	"errors"
	"filmoteka/internal/domain/models"
	"filmoteka/internal/utils"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"strconv"
)

type ActorHandler struct {
	actorUseCase   actorUseCase
	sessionManager *scs.SessionManager
}

func New(useCase actorUseCase, manager *scs.SessionManager) *ActorHandler {
	return &ActorHandler{
		actorUseCase:   useCase,
		sessionManager: manager,
	}
}

type actorUseCase interface {
	GetAllActors() ([]*models.Actor, error)
	CreateActor(a *models.Actor) (*models.Actor, error)
	GetActorByID(id int) (*models.Actor, error)
	UpdateActor(a *models.Actor) (*models.Actor, error)
	DeleteActor(id int) error
}

func (h *ActorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	role := h.sessionManager.GetString(r.Context(), "role")
	if role != "admin" && r.Method != http.MethodGet {
		utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == http.MethodGet && len(r.URL.Query()) > 0:
		h.getActorsById(w, r)
	case r.Method == http.MethodGet:
		h.getAllActors(w, r)
	case r.Method == http.MethodPost:
		h.createActor(w, r)
	case r.Method == http.MethodPatch:
		h.updateActor(w, r)
	case r.Method == http.MethodDelete:
		h.deleteActor(w, r)
	default:
		utils.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	}
}

func (h *ActorHandler) getActorsById(w http.ResponseWriter, r *http.Request) {
	expectedParams := map[string]bool{
		"id": true,
	}

	for param := range r.URL.Query() {
		if !expectedParams[param] {
			log.Println("Invalid request parameter: ", param)
			utils.ErrorJSON(w, errors.New("invalid request parameter "), http.StatusBadRequest)
			return
		}
	}

	var err error
	//utils.ReadJSON(r, w, &actor)
	actorID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Error getting actor, invalid value", err)
		utils.ErrorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}
	actor, err := h.actorUseCase.GetActorByID(actorID)
	if err != nil {
		log.Println("Error getting actor", err)
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Actor retrieved", Data: actor})
}

func (h *ActorHandler) createActor(w http.ResponseWriter, r *http.Request) {
	// Handle post request ==> add new actor
	var actor *models.Actor

	err := utils.ReadJSON(r, w, &actor)
	log.Println("Actor: ", actor)
	if err != nil {
		log.Println("Error reading request", err)
		utils.ErrorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
		return
	}

	_, err = h.actorUseCase.CreateActor(actor)
	if err != nil {
		log.Println("Error creating actor", err)
		utils.ErrorJSON(w, errors.New("error creating actor"), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.JsonResponse{Error: false, Message: "Actor created", Data: actor})
}

func (h *ActorHandler) getAllActors(w http.ResponseWriter, r *http.Request) {
	// Handle get request ==> return all actors
	actors, err := h.actorUseCase.GetAllActors()
	if errors.Is(err, models.ErrNoRecord) {
		utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error getting actors", err)
		utils.ErrorJSON(w, errors.New("Error getting actors"), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Actors retrieved", Data: actors})
}

func (h *ActorHandler) updateActor(w http.ResponseWriter, r *http.Request) {
	actor := &models.Actor{}
	err := utils.ReadJSON(r, w, &actor)
	if err != nil {
		log.Println("Error reading request", err)
		utils.ErrorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
		return
	}
	res, err := h.actorUseCase.UpdateActor(actor)
	if err != nil {
		log.Println("Error updating actor", err)
		utils.ErrorJSON(w, errors.New("error updating actor"), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Actors updated", Data: res})
}

func (h *ActorHandler) deleteActor(w http.ResponseWriter, r *http.Request) {
	expectedParams := map[string]bool{
		"id": true,
	}

	for param := range r.URL.Query() {
		if !expectedParams[param] {
			log.Println("Invalid request parameter: ", param)
			utils.ErrorJSON(w, errors.New("invalid request parameter. "), http.StatusBadRequest)
			return
		}
	}
	var err error
	//actor := &models.Actor{}
	actorID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Error deleting actor, invalid value", err)
		utils.ErrorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}
	err = h.actorUseCase.DeleteActor(actorID)

	if err != nil {
		log.Println("Error deleting actor", err)
		utils.ErrorJSON(w, errors.New("error deleting actor"), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.JsonResponse{Error: false, Message: "Actor successfully deleted"})
}
