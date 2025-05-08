package db

import (
	"Factory/pkg/models"
)

// AddIntegrationCode добавляет новый код интеграции в базу данных
func AddIntegrationCode(code models.IntegrationCode) error {
	query := `
		INSERT INTO integration_codes (
			integration_file_id, code, status, position
		) VALUES (?, ?, ?, ?)
	`

	_, err := DB.Exec(
		query,
		code.IntegrationFileID,
		code.Code,
		code.Status,
		code.Position,
	)

	return err
}

// AddIntegrationCodes добавляет список кодов интеграции в базу данных
func AddIntegrationCodes(fileID int, codes []string) (int, error) {
	// Начинаем транзакцию для повышения производительности
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}

	// Подготавливаем запрос для повторного использования
	stmt, err := tx.Prepare(`
		INSERT INTO integration_codes (
			integration_file_id, code, status, position
		) VALUES (?, ?, ?, ?)
	`)
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

		_, err := stmt.Exec(fileID, code, models.IntegrationCodeStatusPending, position+1)
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

// GetIntegrationCodesByFileID получает коды интеграции по ID файла
func GetIntegrationCodesByFileID(fileID int) ([]models.IntegrationCode, error) {
	query := `
		SELECT id, integration_file_id, code, task_id, status, position
		FROM integration_codes
		WHERE integration_file_id = ?
		ORDER BY position
	`

	rows, err := DB.Query(query, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []models.IntegrationCode

	for rows.Next() {
		var code models.IntegrationCode
		err := rows.Scan(
			&code.ID,
			&code.IntegrationFileID,
			&code.Code,
			&code.TaskID,
			&code.Status,
			&code.Position,
		)
		if err != nil {
			return nil, err
		}

		codes = append(codes, code)
	}

	return codes, nil
}

// GetPendingIntegrationCodes получает неназначенные коды интеграции по ID файла
func GetPendingIntegrationCodes(fileID int) ([]models.IntegrationCode, error) {
	query := `
		SELECT id, integration_file_id, code, task_id, status, position
		FROM integration_codes
		WHERE integration_file_id = ? AND status = ?
		ORDER BY position
	`

	rows, err := DB.Query(query, fileID, models.IntegrationCodeStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []models.IntegrationCode

	for rows.Next() {
		var code models.IntegrationCode
		err := rows.Scan(
			&code.ID,
			&code.IntegrationFileID,
			&code.Code,
			&code.TaskID,
			&code.Status,
			&code.Position,
		)
		if err != nil {
			return nil, err
		}

		codes = append(codes, code)
	}

	return codes, nil
}

// AssignCodesToTask назначает коды интеграции заданию
func AssignCodesToTask(fileID int, taskID int) error {
	query := `
		UPDATE integration_codes
		SET task_id = ?, status = ?
		WHERE integration_file_id = ? AND status = ?
	`

	_, err := DB.Exec(
		query,
		taskID,
		models.IntegrationCodeStatusAssigned,
		fileID,
		models.IntegrationCodeStatusPending,
	)

	return err
}

// CopyCodesFromIntegrationToMarkCodes копирует коды из integration_codes в mark_codes одним запросом
func CopyCodesFromIntegrationToMarkCodes(taskID int) (int, error) {
	// Используем один SQL-запрос для копирования всех кодов
	query := `
        INSERT INTO mark_codes (task_id, code, status, file_position)
        SELECT ?, code, ?, position
        FROM integration_codes
        WHERE task_id = ? AND status = ?
    `

	result, err := DB.Exec(query, taskID, models.MarkCodeStatusNew, taskID, models.IntegrationCodeStatusAssigned)
	if err != nil {
		return 0, err
	}

	// Получаем количество добавленных строк
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(affected), nil
}

// UpdateIntegrationCodeStatus обновляет статус кода интеграции
func UpdateIntegrationCodeStatus(id int, status string) error {
	query := "UPDATE integration_codes SET status = ? WHERE id = ?"
	_, err := DB.Exec(query, status, id)
	return err
}

// CountIntegrationCodesByFileID подсчитывает количество кодов для файла
func CountIntegrationCodesByFileID(fileID int) (int, error) {
	query := "SELECT COUNT(*) FROM integration_codes WHERE integration_file_id = ?"

	var count int
	err := DB.QueryRow(query, fileID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
