package handlers

import (
	"Factory/pkg/db"
	"Factory/pkg/models"
	"Factory/templates"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// MarkCodesListHandler показывает список кодов маркировки для задания
func MarkCodesListHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Неверный ID задания: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем список кодов маркировки
	codes, err := db.GetMarkCodesByTaskID(taskID)
	if err != nil {
		http.Error(w, "Ошибка при получении списка кодов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отображаем шаблон
	if r.Header.Get("HX-Request") == "true" {
		templates.MarkCodesList(codes, taskID).Render(r.Context(), w)
	} else {
		templates.Page(templates.MarkCodesList(codes, taskID)).Render(r.Context(), w)
	}
}

// UploadCodesFormHandler показывает форму для загрузки дополнительных кодов
func UploadCodesFormHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Неверный ID задания: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем информацию о задании
	task, err := db.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "Задание не найдено: "+err.Error(), http.StatusNotFound)
		return
	}

	// Отображаем форму загрузки
	if r.Header.Get("HX-Request") == "true" {
		templates.UploadCodesForm(task).Render(r.Context(), w)
	} else {
		templates.Page(templates.UploadCodesForm(task)).Render(r.Context(), w)
	}
}

// UploadCodesHandler обрабатывает загрузку дополнительных кодов для задания
func UploadCodesHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Неверный ID задания: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Обрабатываем форму с файлом
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Ошибка при обработке формы: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем загруженный файл
	file, handler, err := r.FormFile("mark_codes")
	if err != nil || handler == nil {
		http.Error(w, "Файл не был загружен: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Обрабатываем файл с кодами
	_, err = processMarkCodesFile(taskID, file)
	if err != nil {
		http.Error(w, "Ошибка при обработке файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на страницу с кодами
	http.Redirect(w, r, "/tasks/"+taskIDStr+"/mark-codes", http.StatusSeeOther)
}

// MarkCodeAsUsedHandler отмечает код как использованный
func MarkCodeAsUsedHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания и код из URL
	taskIDStr := chi.URLParam(r, "taskID")
	codeIDStr := chi.URLParam(r, "codeID")

	_, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Неверный ID задания", http.StatusBadRequest)
		return
	}

	codeID, err := strconv.Atoi(codeIDStr)
	if err != nil {
		http.Error(w, "Неверный ID кода", http.StatusBadRequest)
		return
	}

	// Обновляем статус кода
	err = db.UpdateMarkCodeStatus(codeID, models.MarkCodeStatusUsed)
	if err != nil {
		http.Error(w, "Ошибка при обновлении статуса: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на страницу с кодами
	http.Redirect(w, r, "/tasks/"+taskIDStr+"/mark-codes", http.StatusSeeOther)
}

// MarkCodeAsInvalidHandler отмечает код как недействительный
func MarkCodeAsInvalidHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания и код из URL
	taskIDStr := chi.URLParam(r, "taskID")
	codeIDStr := chi.URLParam(r, "codeID")

	_, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Неверный ID задания", http.StatusBadRequest)
		return
	}

	codeID, err := strconv.Atoi(codeIDStr)
	if err != nil {
		http.Error(w, "Неверный ID кода", http.StatusBadRequest)
		return
	}

	// Обновляем статус кода
	err = db.UpdateMarkCodeStatus(codeID, models.MarkCodeStatusInvalid)
	if err != nil {
		http.Error(w, "Ошибка при обновлении статуса: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на страницу с кодами
	http.Redirect(w, r, "/tasks/"+taskIDStr+"/mark-codes", http.StatusSeeOther)
}

// ResetMarkCodeStatusHandler сбрасывает статус кода на "new"
func ResetMarkCodeStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания и код из URL
	taskIDStr := chi.URLParam(r, "taskID")
	codeIDStr := chi.URLParam(r, "codeID")

	_, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Неверный ID задания", http.StatusBadRequest)
		return
	}

	codeID, err := strconv.Atoi(codeIDStr)
	if err != nil {
		http.Error(w, "Неверный ID кода", http.StatusBadRequest)
		return
	}

	// Обновляем статус кода
	err = db.UpdateMarkCodeStatus(codeID, models.MarkCodeStatusNew)
	if err != nil {
		http.Error(w, "Ошибка при обновлении статуса: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на страницу с кодами
	http.Redirect(w, r, "/tasks/"+taskIDStr+"/mark-codes", http.StatusSeeOther)
}
