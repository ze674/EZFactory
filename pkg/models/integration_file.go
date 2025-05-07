package models

import (
	"time"
)

// Константы для статусов файлов интеграции
const (
	FileStatusNew        = "new"        // Новый файл, еще не обработан
	FileStatusProcessing = "processing" // В процессе обработки (задание создано)
	FileStatusCompleted  = "completed"  // Обработка завершена, результаты отправлены
	FileStatusError      = "error"      // Ошибка обработки
)

type IntegrationFile struct {
	ID           int       // Уникальный идентификатор
	UUID         string    // UUID из document_id файла
	Filename     string    // Имя файла
	FilePath     string    // Путь к файлу
	GTIN         string    // GTIN товара из файла
	ProductID    int       // ID соответствующего продукта в системе
	BatchNumber  string    // Номер партии из файла
	Date         string    // Дата производства из файла
	CodesCount   int       // Количество кодов в файле
	Status       string    // Статус обработки файла
	ErrorMessage string    // Сообщение об ошибке (если есть)
	TaskID       int       // ID связанного задания (может быть NULL)
	CreatedAt    time.Time // Дата/время добавления файла в систему
	ProcessedAt  time.Time // Дата/время обработки файла
}
