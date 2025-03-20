package handlers

import (
	"Factory/pkg/db"
	"Factory/templates"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// TaskDetailsHandler обрабатывает запрос на просмотр деталей задания
func TaskDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задания: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем информацию о задании
	task, err := db.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Ошибка при получении информации о задании: "+err.Error(), http.StatusNotFound)
		return
	}

	// Получаем историю изменений
	history, err := db.GetTaskHistory(id)
	if err != nil {
		http.Error(w, "Ошибка при получении истории изменений задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Рендерим шаблон с деталями задания и историей
	templates.TaskDetails(task, history).Render(r.Context(), w)
}
