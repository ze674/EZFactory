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

	// Получаем список производственных линий
	lines, err := db.GetProductionLines()
	if err != nil {
		http.Error(w, "Ошибка при получении списка производственных линий: "+err.Error(), http.StatusInternalServerError)
		return
	}

	templates.AddTaskForm(product, lines).Render(r.Context(), w)
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

	// Получаем ID линии из формы
	lineID, err := strconv.Atoi(r.FormValue("line_id"))
	if err != nil {
		http.Error(w, "Неверный ID производственной линии: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация даты (простая проверка формата)
	if len(date) != 10 {
		http.Error(w, "Неверный формат даты. Используйте формат ДД.ММ.ГГГГ", http.StatusBadRequest)
		return
	}

	task := models.Task{
		ProductID:   id,
		LineID:      lineID,
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
	// Получаем все задания (теперь с информацией о линии)
	tasks, err := db.GetTasks()
	if err != nil {
		http.Error(w, "Ошибка при получении списка заданий: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Рендерим шаблон с заданиями
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

// UpdateTaskStatusWebHandler обрабатывает запрос на изменение статуса задания из веб-интерфейса
func UpdateTaskStatusWebHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задания: "+err.Error(), http.StatusBadRequest)
		return
	}

	newStatus := r.URL.Query().Get("status")
	if newStatus == "" {
		http.Error(w, "Не указан новый статус", http.StatusBadRequest)
		return
	}

	// Обновляем статус в БД (эта функция также сохраняет историю)
	err = db.UpdateTaskStatus(id, newStatus)
	if err != nil {
		http.Error(w, "Ошибка при обновлении статуса задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем обновленные данные о задании
	task, err := db.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Ошибка при получении информации о задании: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем обновленную историю изменений
	history, err := db.GetTaskHistory(id)
	if err != nil {
		http.Error(w, "Ошибка при получении истории изменений задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем обновленный шаблон через HTMX
	templates.TaskDetails(task, history).Render(r.Context(), w)
}
