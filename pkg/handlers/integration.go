package handlers

import (
	"Factory/pkg/db"
	"Factory/pkg/integration"
	"Factory/pkg/models"
	"Factory/templates"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// IntegrationFilesListHandler обрабатывает запрос на просмотр списка входящих файлов из 1С
func IntegrationFilesListHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем список файлов
	files, err := db.GetIntegrationFiles()
	if err != nil {
		http.Error(w, "Ошибка получения списка файлов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отображаем шаблон
	if r.Header.Get("HX-Request") == "true" {
		templates.IntegrationFilesList(files).Render(r.Context(), w)
	} else {
		templates.Page(templates.IntegrationFilesList(files)).Render(r.Context(), w)
	}
}

// IntegrationFileDetailsHandler обрабатывает запрос на просмотр деталей файла
func IntegrationFileDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID файла из URL
	fileIDStr := chi.URLParam(r, "id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, "Неверный ID файла: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле
	file, err := db.GetIntegrationFileByID(fileID)
	if err != nil {
		http.Error(w, "Ошибка получения информации о файле: "+err.Error(), http.StatusNotFound)
		return
	}

	// Получаем информацию о продукте (если найден)
	var product models.Product
	if file.ProductID > 0 {
		product, err = db.GetProductByID(file.ProductID)
		if err != nil {
			// Если продукт не найден, игнорируем ошибку
			fmt.Printf("Продукт с ID %d не найден: %v\n", file.ProductID, err)
		}
	}

	// Получаем список производственных линий для формы создания задания
	lines, err := db.GetProductionLines()
	if err != nil {
		http.Error(w, "Ошибка получения списка производственных линий: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Важное изменение: всегда оборачиваем контент в шаблон Page
	if r.Header.Get("HX-Request") == "true" {
		// Даже для HTMX-запроса возвращаем полную страницу
		templates.Page(templates.IntegrationFileDetails(file, product, []models.IntegrationCode{}, lines)).Render(r.Context(), w)
	} else {
		// Для обычного запроса - также возвращаем полную страницу
		templates.Page(templates.IntegrationFileDetails(file, product, []models.IntegrationCode{}, lines)).Render(r.Context(), w)
	}
}

// CreateTaskFromIntegrationFileHandler обрабатывает запрос на создание задания из файла интеграции
func CreateTaskFromIntegrationFileHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID файла из URL
	fileIDStr := chi.URLParam(r, "id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, "Неверный ID файла: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле
	file, err := db.GetIntegrationFileByID(fileID)
	if err != nil {
		http.Error(w, "Ошибка получения информации о файле: "+err.Error(), http.StatusNotFound)
		return
	}

	// Проверяем, что файл имеет статус "новый"
	if file.Status != models.FileStatusNew {
		http.Error(w, "Невозможно создать задание для файла со статусом "+file.Status, http.StatusBadRequest)
		return
	}

	// Разбираем форму
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка при обработке формы: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем данные из формы
	lineIDStr := r.FormValue("line_id")
	lineID, err := strconv.Atoi(lineIDStr)
	if err != nil {
		http.Error(w, "Неверный ID производственной линии: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, есть ли продукт в системе
	var productID int
	if file.ProductID > 0 {
		productID = file.ProductID
	} else {
		// Если продукт не найден, создаем новый
		product := models.Product{
			Name:      "Продукт из 1С", // Базовое название
			GTIN:      file.GTIN,
			LabelData: "", // Пустые данные для этикетки
		}

		err = db.AddProduct(product)
		if err != nil {
			http.Error(w, "Ошибка при создании продукта: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Получаем ID нового продукта
		newProduct, err := db.GetProductByGTIN(file.GTIN)
		if err != nil {
			http.Error(w, "Ошибка при получении созданного продукта: "+err.Error(), http.StatusInternalServerError)
			return
		}

		productID = newProduct.ID

		// Обновляем запись о файле
		file.ProductID = productID
		db.UpdateIntegrationFileProductID(fileID, productID)
	}

	// Создаем задание
	task := models.Task{
		ProductID:   productID,
		LineID:      lineID,
		Date:        file.Date,
		BatchNumber: file.BatchNumber,
		Status:      "новое",
	}

	// Сохраняем задание в базе данных
	taskID, err := db.AddTask(task)
	if err != nil {
		http.Error(w, "Ошибка при создании задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Связываем файл с заданием
	err = db.LinkFileToTask(fileID, int(taskID))
	if err != nil {
		http.Error(w, "Ошибка при связывании файла с заданием: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Назначение кодов заданию")
	// Назначаем коды заданию
	err = db.AssignCodesToTask(fileID, int(taskID))
	if err != nil {
		http.Error(w, "Ошибка при назначении кодов заданию: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Коды назначены")

	fmt.Println("Копируем коды из таблицы integration_codes в mark_codes")
	// Копируем коды из таблицы integration_codes в mark_codes
	_, err = db.CopyCodesFromIntegrationToMarkCodes(int(taskID))
	if err != nil {
		http.Error(w, "Ошибка при копировании кодов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Коды скопированы")

	// Перенаправляем на страницу задания
	http.Redirect(w, r, fmt.Sprintf("/tasks/%d", taskID), http.StatusSeeOther)
}

// ScanIntegrationDirectoryHandler обрабатывает запрос на сканирование директории с файлами
func ScanIntegrationDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	// Запускаем сканирование директории
	_, err := integration.ScanDirectory(integration.DirectoryPath)
	if err != nil {
		http.Error(w, "Ошибка при сканировании директории: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем список файлов
	files, err := db.GetIntegrationFiles()
	if err != nil {
		http.Error(w, "Ошибка получения списка файлов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем, был ли это HTMX-запрос
	if r.Header.Get("HX-Request") == "true" {
		// Для HTMX-запроса отображаем всю страницу, а не только список
		templates.Page(templates.IntegrationFilesList(files)).Render(r.Context(), w)
	} else {
		// Для обычного запроса - перенаправляем на страницу со списком файлов
		http.Redirect(w, r, "/integration/files", http.StatusSeeOther)
	}
}
