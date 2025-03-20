package handlers

import (
	"Factory/pkg/api"
	"Factory/templates"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SetupRoutes(r *chi.Mux) {
	r.Get("/", homeHandler)
	r.Get("/products", ProductsHandler)
	r.Post("/products", AddProductHandler)
	r.Delete("/products/{id}", DeleteProductHandler)
	r.Get("/products/search", SearchHandler)
	r.Get("/products/{id}/label", LabelHandler)
	r.Post("/products/{id}/label", UpdateLabelHandler)
	r.Get("/products/add-form", AddFormHandler)
	r.Get("/empty", EmptyHandler)
	//Маршруты для заданий
	r.Get("/tasks", TasksListHandler)          // обработчик для просмотра списка заданий
	r.Get("/tasks/{id}", TaskDetailsHandler)   // обработчик для просмотра деталей задания
	r.Delete("/tasks/{id}", DeleteTaskHandler) // обработчик для удаления задания
	r.Get("/products/{id}/add-task", ShowAddTaskForm)
	r.Post("/products/{id}/add-task", AddTaskHandler)
	r.Post("/tasks/{id}/status", api.UpdateTaskStatusHandler)

	//Маршруты для API
	r.Route("/api", func(r chi.Router) {
		r.Get("/tasks", api.GetTasksHandler)
		r.Get("/product/{id}", api.GetProductByIDHandler)
		r.Get("/lines", api.GetProductionLinesHandler)
	})
}

// Здесь можно определить функции-обработчики, если они универсальны
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		templates.Home().Render(r.Context(), w)
	} else {
		templates.Page(templates.Home()).Render(r.Context(), w)
	}
}
