package userhandlers

import (
	"filmoteka/internal/utils"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserHandler struct {
	userUseCase    userUseCase
	sessionManager *scs.SessionManager
}

func New(useCase userUseCase, manager *scs.SessionManager) *UserHandler {
	return &UserHandler{
		userUseCase:    useCase,
		sessionManager: manager,
	}
}

type userUseCase interface {
	//auth(u *models.User)
	GetUserByEmail(email string) (string, string, error)
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle auth here
	switch r.Method {
	case http.MethodPost:
		h.login(w, r)
	default:
		utils.ErrorJSON(w, fmt.Errorf("error method not allowed"), http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}
	err := utils.ReadJSON(r, w, &req)
	log.Println(req)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	pwd, role, err := h.userUseCase.GetUserByEmail(req.Email)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(req.Password))
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// create session
	err = h.sessionManager.RenewToken(r.Context())
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	h.sessionManager.Put(r.Context(), "role", role)

	utils.WriteJSON(w, http.StatusOK, utils.JsonResponse{Error: false, Message: "Logged in", Data: nil})
}
