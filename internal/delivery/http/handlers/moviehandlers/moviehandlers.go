package moviehandlers

import (
	"errors"
	"filmoteka/internal/domain/models"
	"filmoteka/internal/utils"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
)

type MovieHandler struct {
	movieUseCase   movieUseCase
	sessionManager *scs.SessionManager
}

func New(useCase movieUseCase, manager *scs.SessionManager) *MovieHandler {
	return &MovieHandler{
		movieUseCase:   useCase,
		sessionManager: manager,
	}
}

type movieUseCase interface {
	GetAllMovies(param string) ([]*models.Movie, error)
	CreateMovie(m *models.Movie) (*models.Movie, error)
	GetMovieByID(id int) (*models.Movie, error)
	UpdateMovie(m *models.Movie) (*models.Movie, error)
	DeleteMovie(id int) error
	GetMovieByMovieName(moviename string) ([]*models.Movie, error)
}

func (h *MovieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	role := h.sessionManager.GetString(r.Context(), "role")
	//if role == "user" && r.Method != http.MethodGet ||
	if role != "admin" && r.Method != http.MethodGet {
		utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == http.MethodGet:
		h.getMovie(w, r)
	case r.Method == http.MethodPost:
		h.createMovie(w, r)
	case r.Method == http.MethodPatch:
		h.updateMovie(w, r)
	case r.Method == http.MethodDelete:
		h.deleteMovie(w, r)
	default:
		utils.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	}
}

//func (h *MovieHandler) HandleMovies(w http.ResponseWriter, r *http.Request) {
//	role := h.sessionManager.GetString(r.Context(), "role")
//	if role == "user" && r.Method != http.MethodGet ||
//		role != "admin" {
//		utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
//		return
//	}
//
//	switch {
//	case r.Method == http.MethodGet:
//		h.getMovie(w, r)
//	case r.Method == http.MethodPost:
//		h.createMovie(w, r)
//	case r.Method == http.MethodPatch:
//		h.updateMovie(w, r)
//	case r.Method == http.MethodDelete:
//		h.deleteMovie(w, r)
//	default:
//		utils.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
//	}
//}

func (h *MovieHandler) getMovie(w http.ResponseWriter, r *http.Request) {
	expectedParams := map[string]bool{
		"id":        true,
		"sort":      true,
		"name":      true,
		"moviename": true,
	}

	for param := range r.URL.Query() {
		if !expectedParams[param] {
			log.Println("Invalid request parameter: ", param)
			utils.ErrorJSON(w, errors.New("invalid request parameter. "), http.StatusBadRequest)
			return
		}
	}
	idParam, idOk := r.URL.Query()["id"]
	sortParam, sortOk := r.URL.Query()["sort"]
	nameParam, nameOk := r.URL.Query()["name"]

	if (sortOk && idOk) || (sortOk && nameOk) || (nameOk && idOk) || (sortOk && idOk && nameOk) {
		utils.ErrorJSON(w, errors.New("invalid request"), http.StatusBadRequest)

	} else if idOk {
		// Fetch movie by id
		var id int
		utils.StringToInt(w, &id, idParam[0])
		movie, err := h.movieUseCase.GetMovieByID(id)
		if err != nil {
			log.Println("Error getting movie", err)
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Movie retrieved", Data: movie})

	} else if sortOk {
		// Fetch all movies and sort
		movies, err := h.movieUseCase.GetAllMovies(sortParam[0])
		if errors.Is(err, models.ErrNoRecord) {
			utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
			return
		} else if err != nil {
			log.Println("Error getting movies", err)
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Movies retrieved", Data: movies})

	} else if nameOk && len(nameParam) > 0 {
		// Fetch all movies and sort
		movies, err := h.movieUseCase.GetMovieByMovieName(nameParam[0])
		if errors.Is(err, models.ErrNoRecord) {
			utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
			return
		} else if err != nil {
			log.Println("Error getting movies", err)
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Movie retrieved", Data: movies})
	} else {
		// Fetch all movies
		movies, err := h.movieUseCase.GetAllMovies("")
		if errors.Is(err, models.ErrNoRecord) {
			utils.ErrorJSON(w, models.ErrNoRecord, http.StatusNotFound)
			return
		} else if err != nil {
			log.Println("Error getting movies", err)
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Movies retrieved", Data: movies})
	}
}

func (h *MovieHandler) createMovie(w http.ResponseWriter, r *http.Request) {
	// Handle post request ==> add new actor
	var movie *models.Movie

	err := utils.ReadJSON(r, w, &movie)
	log.Println("Movie: ", movie)
	if err != nil {
		log.Println("Error reading request", err)
		utils.ErrorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
		return
	}

	movie, err = h.movieUseCase.CreateMovie(movie)
	if err != nil {
		log.Println("Error creating movie", err)
		utils.ErrorJSON(w, errors.New("error creating movie"), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.JsonResponse{Error: false, Message: "Movie created", Data: movie})

}

// handler for a patch request
func (h *MovieHandler) updateMovie(w http.ResponseWriter, r *http.Request) {
	movie := &models.Movie{}
	err := utils.ReadJSON(r, w, &movie)
	if err != nil {
		log.Println("Error reading request", err)
		utils.ErrorJSON(w, errors.New("error reading request"), http.StatusBadRequest)
		return
	}
	movie, err = h.movieUseCase.UpdateMovie(movie)
	if err != nil {
		log.Println("Error updating movie", err)
		utils.ErrorJSON(w, errors.New("error updating movie"), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Movie updated", Data: movie})
}

func (h *MovieHandler) deleteMovie(w http.ResponseWriter, r *http.Request) {
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
	idParam := r.URL.Query()["id"]
	if len(idParam) < 1 {
		utils.ErrorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return

	}
	var err error
	var id int
	utils.StringToInt(w, &id, idParam[0])
	err = h.movieUseCase.DeleteMovie(id)
	if err != nil {
		log.Println("Error deleting movie", err)
		utils.ErrorJSON(w, errors.New("error deleting movie"), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.JsonResponse{Error: false, Message: "Movie successfully deleted"})
}
