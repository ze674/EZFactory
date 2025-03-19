package handlers

import (
	"net/http"
	"strconv"
	"time"

	"Factory/pkg/db"
	"Factory/pkg/models"
	"Factory/templates"
	"github.com/go-chi/chi/v5"
)

// ShowAddTaskForm показывает форму для создания задания
func ShowAddTaskForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID продукта: "+err.Error(), http.StatusBadRequest)
		return
	}

	product, err := db.GetProductByID(id)
	if err != nil {
		http.Error(w, "Ошибка при получении продукта: "+err.Error(), http.StatusNotFound)
		return
	}

	templates.AddTaskForm(product).Render(r.Context(), w)
}

// AddTaskHandler обрабатывает отправку формы создания задания
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID продукта: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка при обработке формы: "+err.Error(), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	batchNumber := r.FormValue("batch_number")
	status := r.FormValue("status")

	// Валидация даты (простая проверка формата)
	if len(date) != 10 {
		http.Error(w, "Неверный формат даты. Используйте формат ДД.ММ.ГГГГ", http.StatusBadRequest)
		return
	}

	task := models.Task{
		ProductID:   id,
		Date:        date,
		BatchNumber: batchNumber,
		Status:      status,
		CreatedAt:   time.Now(),
	}

	err = db.AddTask(task)
	if err != nil {
		http.Error(w, "Ошибка при создании задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на список продуктов
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

// TasksListHandler обрабатывает запрос на просмотр списка заданий
func TasksListHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks()
	if err != nil {
		http.Error(w, "Ошибка при получении списка заданий: "+err.Error(), http.StatusInternalServerError)
		return
	}

	templates.TasksList(tasks).Render(r.Context(), w)
}

// DeleteTaskHandler обрабатывает запрос на удаление задания
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задания: "+err.Error(), http.StatusBadRequest)
		return
	}

	success, err := db.DeleteTask(id)
	if err != nil {
		http.Error(w, "Ошибка при удалении задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !success {
		http.Error(w, "Задание не найдено", http.StatusNotFound)
		return
	}

	// Для HTMX запросов отправляем обновленный список заданий
	tasks, err := db.GetTasks()
	if err != nil {
		http.Error(w, "Ошибка при получении списка заданий: "+err.Error(), http.StatusInternalServerError)
		return
	}

	templates.TasksList(tasks).Render(r.Context(), w)

}
