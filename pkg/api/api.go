package api

import (
	"encoding/json"
	"net/http"
)

// Response - стандартная структура ответа API
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RespondWithJSON - вспомогательная функция для отправки JSON-ответа
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"Ошибка при формировании ответа"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError - вспомогательная функция для отправки ошибки
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, Response{
		Success: false,
		Error:   message,
	})
}
