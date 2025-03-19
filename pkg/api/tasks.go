package api

import (
	"Factory/pkg/db"
	"net/http"
)

// GetTasksHandler - обработчик для получения списка заданий через API
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении списка заданий: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    tasks,
	})
}
