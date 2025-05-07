package models

// Константы для статусов кодов интеграции
const (
	IntegrationCodeStatusPending  = "pending"  // Ожидает обработки
	IntegrationCodeStatusAssigned = "assigned" // Назначен заданию
	IntegrationCodeStatusUsed     = "used"     // Использован (отсканирован)
)

type IntegrationCode struct {
	ID                int    // Уникальный идентификатор
	IntegrationFileID int    // ID файла интеграции
	Code              string // Код маркировки
	TaskID            int    // ID задания (может быть 0, если не назначен)
	Status            string // Статус кода
	Position          int    // Позиция в исходном файле
}
