package api

import (
	"Factory/pkg/db"
	"net/http"
)

// GetProductionLinesHandler - обработчик для получения списка производственных линий через API
func GetProductionLinesHandler(w http.ResponseWriter, r *http.Request) {
	lines, err := db.GetProductionLines()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении списка линий: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    lines,
	})
}
