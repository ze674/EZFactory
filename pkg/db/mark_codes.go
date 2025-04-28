package db

import (
	"Factory/pkg/models"
	"database/sql"
	"time"
)

// AddMarkCode добавляет новый код маркировки в базу данных
func AddMarkCode(taskID int, code string, position int) error {
	query := "INSERT INTO mark_codes (task_id, code, status, file_position) VALUES (?, ?, ?, ?)"
	_, err := DB.Exec(query, taskID, code, models.MarkCodeStatusNew, position)
	return err
}

// AddMarkCodes добавляет список кодов маркировки для задания
func AddMarkCodes(taskID int, codes []string) (int, error) {
	// Начинаем транзакцию для повышения производительности
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}

	// Подготавливаем запрос для повторного использования
	stmt, err := tx.Prepare("INSERT INTO mark_codes (task_id, code, status, file_position) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	// Счетчик добавленных кодов
	added := 0

	// Добавляем каждый код
	for position, code := range codes {
		// Пропускаем пустые строки
		if code == "" {
			continue
		}

		_, err := stmt.Exec(taskID, code, models.MarkCodeStatusNew, position+1) // +1 чтобы начинать с 1, а не с 0
		// Если код уже существует, пропускаем его и продолжаем
		if err != nil {
			continue
		}
		added++
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return added, nil
}

// GetMarkCodesByTaskID возвращает все коды маркировки для заданного задания,
// отсортированные по позиции в файле
func GetMarkCodesByTaskID(taskID int) ([]models.MarkCode, error) {
	query := `
		SELECT id, task_id, code, status, used_at, file_position 
		FROM mark_codes 
		WHERE task_id = ? 
		ORDER BY file_position
	`

	rows, err := DB.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []models.MarkCode

	for rows.Next() {
		var code models.MarkCode
		var usedAt sql.NullTime

		err := rows.Scan(&code.ID, &code.TaskID, &code.Code, &code.Status, &usedAt, &code.FilePosition)
		if err != nil {
			return nil, err
		}

		// Если значение usedAt не NULL, устанавливаем его в структуре
		if usedAt.Valid {
			code.UsedAt = usedAt.Time
		}

		codes = append(codes, code)
	}

	return codes, nil
}

// UpdateMarkCodeStatus обновляет статус кода маркировки
func UpdateMarkCodeStatus(id int, status string) error {
	var query string
	var args []interface{}

	// Если устанавливается статус 'used', устанавливаем также время использования
	if status == models.MarkCodeStatusUsed {
		query = "UPDATE mark_codes SET status = ?, used_at = ? WHERE id = ?"
		args = []interface{}{status, time.Now(), id}
	} else {
		query = "UPDATE mark_codes SET status = ? WHERE id = ?"
		args = []interface{}{status, id}
	}

	_, err := DB.Exec(query, args...)
	return err
}

// GetMarkCodeByCode находит код маркировки по его значению
func GetMarkCodeByCode(code string) (*models.MarkCode, error) {
	query := `
		SELECT id, task_id, code, status, used_at, file_position 
		FROM mark_codes 
		WHERE code = ?
	`

	var markCode models.MarkCode
	var usedAt sql.NullTime

	err := DB.QueryRow(query, code).Scan(
		&markCode.ID,
		&markCode.TaskID,
		&markCode.Code,
		&markCode.Status,
		&usedAt,
		&markCode.FilePosition,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Код не найден
		}
		return nil, err
	}

	// Если значение usedAt не NULL, устанавливаем его в структуре
	if usedAt.Valid {
		markCode.UsedAt = usedAt.Time
	}

	return &markCode, nil
}

// CountMarkCodesByTaskAndStatus подсчитывает количество кодов маркировки
// для задания с определенным статусом
func CountMarkCodesByTaskAndStatus(taskID int, status string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM mark_codes 
		WHERE task_id = ? AND status = ?
	`

	var count int
	err := DB.QueryRow(query, taskID, status).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetMarkCodeStats возвращает статистику по кодам для задания
func GetMarkCodeStats(taskID int) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) 
		FROM mark_codes 
		WHERE task_id = ? 
		GROUP BY status
	`

	rows, err := DB.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)

	for rows.Next() {
		var status string
		var count int

		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, err
		}

		stats[status] = count
	}

	return stats, nil
}
