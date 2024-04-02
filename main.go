package main

import (
	"api/messages"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.DefaultLogger)

	r.Route("/messages", messages.Router)

	http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)(r))
}
