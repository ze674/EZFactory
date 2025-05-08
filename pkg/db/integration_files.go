package db

import (
	"Factory/pkg/models"
	"database/sql"
	"time"
)

// AddIntegrationFile добавляет новый файл интеграции в базу данных
func AddIntegrationFile(file models.IntegrationFile) (int64, error) {
	query := `
		INSERT INTO integration_files (
			uuid, filename, file_path, gtin, product_id, batch_number, 
			date, codes_count, status, error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := DB.Exec(
		query,
		file.UUID,
		file.Filename,
		file.FilePath,
		file.GTIN,
		file.ProductID,
		file.BatchNumber,
		file.Date,
		file.CodesCount,
		file.Status,
		file.ErrorMessage,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetIntegrationFileByID получает файл интеграции по ID
func GetIntegrationFileByID(id int) (models.IntegrationFile, error) {
	query := `
		SELECT id, uuid, filename, file_path, gtin, product_id, batch_number, 
		       date, codes_count, status, error_message, task_id, created_at, processed_at
		FROM integration_files
		WHERE id = ?
	`

	var file models.IntegrationFile
	var productID sql.NullInt64
	var taskID sql.NullInt64
	var processedAt sql.NullString
	var createdAt string

	err := DB.QueryRow(query, id).Scan(
		&file.ID, &file.UUID, &file.Filename, &file.FilePath, &file.GTIN,
		&productID, &file.BatchNumber, &file.Date, &file.CodesCount,
		&file.Status, &file.ErrorMessage, &taskID, &createdAt, &processedAt,
	)
	if err != nil {
		return file, err
	}

	// Обработка NULL значений
	if productID.Valid {
		file.ProductID = int(productID.Int64)
	}

	if taskID.Valid {
		file.TaskID = int(taskID.Int64)
	}

	// Преобразование строки времени в time.Time
	file.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		file.CreatedAt = time.Now()
	}

	if processedAt.Valid {
		file.ProcessedAt, err = time.Parse("2006-01-02 15:04:05", processedAt.String)
		if err != nil {
			file.ProcessedAt = time.Time{}
		}
	}

	return file, nil
}

// GetIntegrationFileByUUID получает файл интеграции по UUID
func GetIntegrationFileByUUID(uuid string) (models.IntegrationFile, error) {
	query := `
		SELECT id, uuid, filename, file_path, gtin, product_id, batch_number, 
		       date, codes_count, status, error_message, task_id, created_at, processed_at
		FROM integration_files
		WHERE uuid = ?
	`

	var file models.IntegrationFile
	var productID, taskID, processedAt interface{}
	var createdAt string

	err := DB.QueryRow(query, uuid).Scan(
		&file.ID, &file.UUID, &file.Filename, &file.FilePath, &file.GTIN,
		&productID, &file.BatchNumber, &file.Date, &file.CodesCount,
		&file.Status, &file.ErrorMessage, &taskID, &createdAt, &processedAt,
	)
	if err != nil {
		return file, err
	}

	// Обработка NULL значений
	if productID != nil {
		file.ProductID = productID.(int)
	}

	if taskID != nil {
		file.TaskID = taskID.(int)
	}

	// Преобразование строки времени в time.Time
	file.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		file.CreatedAt = time.Now()
	}

	if processedAt != nil {
		pTime, err := time.Parse("2006-01-02 15:04:05", processedAt.(string))
		if err == nil {
			file.ProcessedAt = pTime
		}
	}

	return file, nil
}

// GetIntegrationFiles получает список всех файлов интеграции
func GetIntegrationFiles() ([]models.IntegrationFile, error) {
	query := `
		SELECT id, uuid, filename, file_path, gtin, product_id, batch_number, 
		       date, codes_count, status, error_message, task_id, created_at, processed_at
		FROM integration_files
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.IntegrationFile

	for rows.Next() {
		var file models.IntegrationFile
		var productID sql.NullInt64
		var taskID sql.NullInt64
		var processedAt sql.NullString
		var createdAt string

		err := rows.Scan(
			&file.ID, &file.UUID, &file.Filename, &file.FilePath, &file.GTIN,
			&productID, &file.BatchNumber, &file.Date, &file.CodesCount,
			&file.Status, &file.ErrorMessage, &taskID, &createdAt, &processedAt,
		)
		if err != nil {
			return nil, err
		}

		// Обработка NULL значений
		if productID.Valid {
			file.ProductID = int(productID.Int64)
		}

		if taskID.Valid {
			file.TaskID = int(taskID.Int64)
		}

		// Преобразование строки времени в time.Time
		file.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			file.CreatedAt = time.Now()
		}

		if processedAt.Valid {
			file.ProcessedAt, err = time.Parse("2006-01-02 15:04:05", processedAt.String)
			if err != nil {
				file.ProcessedAt = time.Time{}
			}
		}

		files = append(files, file)
	}

	return files, nil
}

// GetNewIntegrationFiles получает список новых файлов интеграции (со статусом "new")
func GetNewIntegrationFiles() ([]models.IntegrationFile, error) {
	query := `
		SELECT id, uuid, filename, file_path, gtin, product_id, batch_number, 
		       date, codes_count, status, error_message, task_id, created_at, processed_at
		FROM integration_files
		WHERE status = ?
		ORDER BY created_at ASC
	`

	rows, err := DB.Query(query, models.FileStatusNew)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.IntegrationFile

	for rows.Next() {
		var file models.IntegrationFile
		var productID, taskID, processedAt interface{}
		var createdAt string

		err := rows.Scan(
			&file.ID, &file.UUID, &file.Filename, &file.FilePath, &file.GTIN,
			&productID, &file.BatchNumber, &file.Date, &file.CodesCount,
			&file.Status, &file.ErrorMessage, &taskID, &createdAt, &processedAt,
		)
		if err != nil {
			return nil, err
		}

		// Обработка NULL значений
		if productID != nil {
			file.ProductID = productID.(int)
		}

		if taskID != nil {
			file.TaskID = taskID.(int)
		}

		// Преобразование строки времени в time.Time
		file.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			file.CreatedAt = time.Now()
		}

		if processedAt != nil {
			pTime, err := time.Parse("2006-01-02 15:04:05", processedAt.(string))
			if err == nil {
				file.ProcessedAt = pTime
			}
		}

		files = append(files, file)
	}

	return files, nil
}

// UpdateIntegrationFileStatus обновляет статус файла интеграции
func UpdateIntegrationFileStatus(id int, status string, errorMessage string) error {
	query := `
		UPDATE integration_files 
		SET status = ?, error_message = ?, processed_at = ?
		WHERE id = ?
	`

	_, err := DB.Exec(query, status, errorMessage, time.Now(), id)
	return err
}

// LinkFileToTask связывает файл интеграции с заданием
func LinkFileToTask(fileID int, taskID int) error {
	query := `
		UPDATE integration_files 
		SET task_id = ?, status = ?, processed_at = ?
		WHERE id = ?
	`

	_, err := DB.Exec(query, taskID, models.FileStatusProcessing, time.Now(), fileID)
	return err
}

// FileExists проверяет, существует ли файл с заданным UUID
func FileExists(uuid string) (bool, error) {
	query := "SELECT COUNT(*) FROM integration_files WHERE uuid = ?"

	var count int
	err := DB.QueryRow(query, uuid).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateIntegrationFileProductID обновляет ID продукта для файла интеграции
func UpdateIntegrationFileProductID(fileID int, productID int) error {
	query := "UPDATE integration_files SET product_id = ? WHERE id = ?"
	_, err := DB.Exec(query, productID, fileID)
	return err
}
