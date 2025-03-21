package main

import (
	"Factory/pkg/db"
	"Factory/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.DB.Close()
	r := chi.NewRouter()
	// Обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	handlers.SetupRoutes(r)

	log.Fatal(http.ListenAndServe(":8081", r))
}
