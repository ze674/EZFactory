package main

import (
	"Factory/pkg/db"
	"Factory/pkg/handlers"
	"Factory/pkg/integration"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	dbFilePath := "./factory.db"
	if _, err := os.Stat(dbFilePath); err == nil {
		fmt.Println("Удаление существующей базы данных...")
		err := os.Remove(dbFilePath)
		if err != nil {
			log.Fatalf("Ошибка при удалении файла базы данных: %v", err)
		}
	}
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.DB.Close()

	// Инициализация директории для интеграции
	err := os.MkdirAll(integration.DirectoryPath, 0755)
	if err != nil {
		log.Fatal("Failed to create integration directory: ", err)
	}

	// Запуск сканера для интеграции
	integration.StartScanner()

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
