package handlers

import (
	"Factory/pkg/db"
	"Factory/templates"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

// Обновленная функция TaskDetailsHandler
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

	// Получаем статистику по кодам маркировки
	codeStats, err := db.GetMarkCodeStats(id)
	if err != nil {
		// Если возникла ошибка, просто логируем её и продолжаем без статистики
		log.Printf("Ошибка при получении статистики кодов: %v", err)
		codeStats = nil
	}

	// Рендерим шаблон с деталями задания, историей и статистикой кодов
	if r.Header.Get("HX-Request") == "true" {
		templates.TaskDetails(task, history, codeStats).Render(r.Context(), w)
	} else {
		templates.Page(templates.TaskDetails(task, history, codeStats)).Render(r.Context(), w)
	}
}
