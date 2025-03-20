package db

import (
	"Factory/pkg/models"
)

// GetProductionLines возвращает список всех производственных линий
func GetProductionLines() ([]models.ProductionLine, error) {
	query := "SELECT id, name FROM production_lines ORDER BY name"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lines := []models.ProductionLine{}

	for rows.Next() {
		var line models.ProductionLine
		err := rows.Scan(&line.ID, &line.Name)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

// GetProductionLineByID возвращает производственную линию по ID
func GetProductionLineByID(id int) (models.ProductionLine, error) {
	query := "SELECT id, name FROM production_lines WHERE id = ?"

	row := DB.QueryRow(query, id)

	var line models.ProductionLine
	err := row.Scan(&line.ID, &line.Name)
	if err != nil {
		return line, err
	}

	return line, nil
}

// AddProductionLine добавляет новую производственную линию
func AddProductionLine(line models.ProductionLine) error {
	query := "INSERT INTO production_lines (name) VALUES (?)"

	result, err := DB.Exec(query, line.Name)
	if err != nil {
		return err
	}

	// Проверка успешности добавления
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// DeleteProductionLine удаляет производственную линию по ID
func DeleteProductionLine(id int) (bool, error) {
	query := "DELETE FROM production_lines WHERE id = ?"

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

// EnsureDefaultProductionLines создает несколько линий по умолчанию, если их нет в базе
func EnsureDefaultProductionLines() error {
	// Проверяем, есть ли линии в базе
	lines, err := GetProductionLines()
	if err != nil {
		return err
	}

	// Если линий нет, добавляем несколько по умолчанию
	if len(lines) == 0 {
		defaultLines := []models.ProductionLine{
			{Name: "Линия 1"},
			{Name: "Линия 2"},
			{Name: "Линия 3"},
		}

		for _, line := range defaultLines {
			err := AddProductionLine(line)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
