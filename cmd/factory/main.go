package main

import (
	"Factory/pkg/db"
	"Factory/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.DB.Close()
	r := chi.NewRouter()
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	handlers.SetupRoutes(r)

	log.Fatal(http.ListenAndServe(":8081", r))
}
