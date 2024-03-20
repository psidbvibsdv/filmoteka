package routes

import (
	"filmoteka/internal/delivery/http/handlers/actorhandlers"
	"filmoteka/internal/delivery/http/handlers/actormoviehandlers"
	"filmoteka/internal/delivery/http/handlers/moviehandlers"
	"filmoteka/internal/delivery/http/handlers/userhandlers"
	"filmoteka/internal/domain/usecase"
	"github.com/alexedwards/scs/v2"
	"net/http"
)

func Routes(useCase *usecase.UseCase, manager *scs.SessionManager) http.Handler {
	mux := http.NewServeMux()
	userHandler := userhandlers.New(useCase.UserUseCase, manager)
	mux.Handle("/login", userHandler)

	actormovieHandler := actormoviehandlers.New(useCase.ActorMovieUseCase, manager)
	mux.Handle("/movie/actormovie", actormovieHandler)

	actorHandler := actorhandlers.New(useCase.ActorUseCase, manager)
	mux.Handle("/actor", actorHandler)

	movieHandler := moviehandlers.New(useCase.MovieUseCase, manager)
	mux.Handle("/movie", movieHandler)

	return manager.LoadAndSave(mux)
}
