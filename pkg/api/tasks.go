package api

import (
	"Factory/pkg/db"
	"Factory/pkg/models"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// GetTasksHandler - обработчик для получения списка заданий через API
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, есть ли параметр line_id в запросе
	lineIDStr := r.URL.Query().Get("line_id")

	var tasks []models.Task
	var err error

	// Если параметр line_id указан, фильтруем задания по линии
	if lineIDStr != "" {
		lineID, err := strconv.Atoi(lineIDStr)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Некорректный ID линии: "+err.Error())
			return
		}

		tasks, err = db.GetTasksByLineID(lineID)
	} else {
		// Иначе получаем все задания
		tasks, err = db.GetTasks()
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении списка заданий: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    tasks,
	})
}

// UpdateTaskStatusHandler обновляет статус задания
func UpdateTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Некорректный ID задания")
		return
	}

	// Получаем новый статус из запроса
	if err := r.ParseForm(); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Не удалось обработать данные запроса")
		return
	}

	newStatus := r.FormValue("status")
	if newStatus == "" {
		RespondWithError(w, http.StatusBadRequest, "Не указан новый статус")
		return
	}

	// Проверяем, что статус допустимый
	if newStatus != "новое" && newStatus != "в работе" && newStatus != "завершено" {
		RespondWithError(w, http.StatusBadRequest, "Недопустимый статус. Разрешенные статусы: новое, в работе, завершено")
		return
	}

	// Обновляем статус в БД
	err = db.UpdateTaskStatus(taskID, newStatus)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Ошибка при обновлении статуса задания: "+err.Error())
		return
	}

	// Возвращаем успешный ответ
	RespondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "Статус задания успешно обновлен",
	})
}

// GetTaskByIDHandler - обработчик для получения информации о задании по ID через API
func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания из URL параметров
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Некорректный ID задания")
		fmt.Println(err)
		return
	}

	// Получаем информацию о задании из базы данных
	task, err := db.GetTaskByID(taskID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении информации о задании: "+err.Error())
		fmt.Println(err)
		return
	}

	// Отправляем успешный ответ
	RespondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    task,
	})
}
