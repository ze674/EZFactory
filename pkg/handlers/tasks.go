package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
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

// pkg/handlers/tasks.go - обновленная функция AddTaskHandler
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID продукта: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Используем ParseMultipartForm для поддержки загрузки файлов
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Ограничение размера файла ~10MB
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

	// Создаем задание в базе данных
	taskID, err := db.AddTask(task)
	if err != nil {
		http.Error(w, "Ошибка при создании задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Обработка загруженного файла с кодами маркировки
	file, handler, err := r.FormFile("mark_codes")
	if err == nil && handler != nil { // Файл был загружен
		defer file.Close()

		// Обрабатываем файл с кодами
		codesCount, err := processMarkCodesFile(int(taskID), file)
		if err != nil {
			// Логгируем ошибку, но продолжаем выполнение
			// (задание уже создано, просто не удалось загрузить коды)
			log.Printf("Ошибка при обработке файла с кодами: %v", err)
		} else {
			log.Printf("Успешно загружено %d кодов маркировки для задания #%d", codesCount, taskID)
		}
	}

	// Перенаправляем на список продуктов
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

// Новая функция для обработки файла с кодами маркировки
func processMarkCodesFile(taskID int, file multipart.File) (int, error) {
	// Создаем буферизированный reader для файла
	reader := bufio.NewReader(file)

	// Читаем все строки из файла
	var codes []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Конец файла, добавляем последнюю строку, если она не пустая
				line = strings.TrimSpace(line)
				if line != "" {
					codes = append(codes, line)
				}
				break
			}
			return 0, fmt.Errorf("ошибка чтения файла: %w", err)
		}

		// Очищаем строку от лишних пробелов и символов перевода строки
		line = strings.TrimSpace(line)
		if line != "" {
			codes = append(codes, line)
		}
	}

	// Если в файле нет кодов, сразу возвращаем 0
	if len(codes) == 0 {
		return 0, nil
	}

	// Добавляем коды в базу данных
	codesAdded, err := db.AddMarkCodes(taskID, codes)
	if err != nil {
		return 0, fmt.Errorf("ошибка сохранения кодов в базе данных: %w", err)
	}

	return codesAdded, nil
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
	// Получаем статистику по кодам маркировки
	codeStats, err := db.GetMarkCodeStats(id)
	if err != nil {
		// Если возникла ошибка, просто логируем её и продолжаем без статистики
		log.Printf("Ошибка при получении статистики кодов: %v", err)
		codeStats = nil
	}

	// Возвращаем обновленный шаблон через HTMX
	templates.TaskDetails(task, history, codeStats).Render(r.Context(), w)
}
