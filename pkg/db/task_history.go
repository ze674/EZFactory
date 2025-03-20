package db

import (
	"Factory/pkg/models"
	"time"
)

// AddTaskHistory добавляет запись в историю изменений задания
func AddTaskHistory(taskID int, oldStatus, newStatus string) error {
	query := "INSERT INTO task_history (task_id, old_status, new_status) VALUES (?, ?, ?)"

	_, err := DB.Exec(query, taskID, oldStatus, newStatus)
	if err != nil {
		return err
	}

	return nil
}

// GetTaskHistory возвращает историю изменений задания
func GetTaskHistory(taskID int) ([]models.TaskHistory, error) {
	query := `
        SELECT id, task_id, old_status, new_status, changed_at
        FROM task_history
        WHERE task_id = ?
        ORDER BY changed_at DESC
    `

	rows, err := DB.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []models.TaskHistory{}

	for rows.Next() {
		var h models.TaskHistory
		var changedAt string

		err := rows.Scan(&h.ID, &h.TaskID, &h.OldStatus, &h.NewStatus, &changedAt)
		if err != nil {
			return nil, err
		}

		// Преобразуем строку времени в time.Time
		h.ChangedAt, err = time.Parse("2006-01-02 15:04:05", changedAt)
		if err != nil {
			// Если не удалось распарсить, используем текущее время
			h.ChangedAt = time.Now()
		}

		history = append(history, h)
	}

	return history, nil
}
