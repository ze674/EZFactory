package api

import (
	"Factory/pkg/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// GetProductByIDHandler обрабатывает запросы на получение информации о продукте по ID
func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID продукта из URL параметров
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Некорректный ID продукта")
		return
	}

	// Получаем информацию о продукте из базы данных
	product, err := db.GetProductByID(productID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении информации о продукте: "+err.Error())
		return
	}

	// Отправляем успешный ответ
	RespondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    product,
	})
}
