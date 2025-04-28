package db

import (
	"Factory/pkg/models"
	"fmt"
	"time"
)

func GetTasks() ([]models.Task, error) {
	query := `
        SELECT t.id, t.product_id, p.name, t.line_id, l.name, t.date, t.batch_number, t.status, t.created_at 
        FROM tasks t
        JOIN products p ON t.product_id = p.id
        JOIN production_lines l ON t.line_id = l.id
        ORDER BY t.created_at DESC
    `

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var t models.Task
		var createdAt string

		err := rows.Scan(&t.ID, &t.ProductID, &t.ProductName, &t.LineID, &t.LineName,
			&t.Date, &t.BatchNumber, &t.Status, &createdAt)
		if err != nil {
			return nil, err
		}

		// Преобразуем строку времени в time.Time
		t.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			// Если не удалось распарсить, используем текущее время
			t.CreatedAt = time.Now()
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func GetTasksByLineID(lineID int) ([]models.Task, error) {
	query := `
        SELECT t.id, t.product_id, p.name, t.line_id, l.name, t.date, t.batch_number, t.status, t.created_at 
        FROM tasks t
        JOIN products p ON t.product_id = p.id
        JOIN production_lines l ON t.line_id = l.id
        WHERE t.line_id = ?
        ORDER BY t.created_at DESC
    `

	rows, err := DB.Query(query, lineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var t models.Task
		var createdAt string

		err := rows.Scan(&t.ID, &t.ProductID, &t.ProductName, &t.LineID, &t.LineName,
			&t.Date, &t.BatchNumber, &t.Status, &createdAt)
		if err != nil {
			return nil, err
		}

		// Преобразуем строку времени в time.Time
		t.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			// Если не удалось распарсить, используем текущее время
			t.CreatedAt = time.Now()
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

// pkg/db/tasks.go - обновленная функция AddTask
func AddTask(task models.Task) (int64, error) {
	query := "INSERT INTO tasks (product_id, line_id, date, batch_number, status) VALUES (?, ?, ?, ?, ?)"

	result, err := DB.Exec(query, task.ProductID, task.LineID, task.Date, task.BatchNumber, task.Status)
	if err != nil {
		return 0, err
	}

	// Получаем ID вставленной записи
	taskID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

func DeleteTask(id int) (bool, error) {
	query := "DELETE FROM tasks WHERE id = ?"

	result, err := DB.Exec(query, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// UpdateTaskStatus обновляет статус задания по ID и сохраняет историю изменений
func UpdateTaskStatus(taskID int, newStatus string) error {
	// Сначала получаем текущий статус задания
	task, err := GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// Если статус не изменился, ничего не делаем
	if task.Status == newStatus {
		return nil
	}

	// Обновляем статус
	query := "UPDATE tasks SET status = ? WHERE id = ?"

	result, err := DB.Exec(query, newStatus, taskID)
	if err != nil {
		return err
	}

	// Проверяем, было ли обновлено задание
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задание с ID %d не найдено", taskID)
	}

	// Добавляем запись в историю изменений
	err = AddTaskHistory(taskID, task.Status, newStatus)
	if err != nil {
		return err
	}

	return nil
}

// GetTaskByID возвращает задание по ID
func GetTaskByID(taskID int) (models.Task, error) {
	query := `
        SELECT t.id, t.product_id, p.name, t.line_id, l.name, t.date, t.batch_number, t.status, t.created_at 
        FROM tasks t
        JOIN products p ON t.product_id = p.id
        JOIN production_lines l ON t.line_id = l.id
        WHERE t.id = ?
    `

	row := DB.QueryRow(query, taskID)

	var t models.Task
	var createdAt string

	err := row.Scan(&t.ID, &t.ProductID, &t.ProductName, &t.LineID, &t.LineName,
		&t.Date, &t.BatchNumber, &t.Status, &createdAt)
	if err != nil {
		return t, err
	}

	// Преобразуем строку времени в time.Time
	t.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		// Если не удалось распарсить, используем текущее время
		t.CreatedAt = time.Now()
	}

	return t, nil
}
