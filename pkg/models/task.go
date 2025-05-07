package models

import (
	"time"
)

const (
	TaskStatusNew        = "новое"
	TaskStatusInProgress = "в работе"
	TaskStatusPaused     = "приостановлено"
	TaskStatusCompleted  = "завершено"
	TaskStatusSent       = "отправлено"
)

type Task struct {
	ID          int       // Уникальный идентификатор
	ProductID   int       // Связь с продуктом
	ProductName string    // Название продукта (для удобства отображения)
	LineID      int       // Связь с производственной линией
	LineName    string    // Название линии (для удобства отображения)
	Date        string    // Дата в формате ДД.ММ.ГГГГ
	BatchNumber string    // Номер партии
	Status      string    // Статус: "новое", "в работе", "завершено", "отправлено", "приостановлено"
	CreatedAt   time.Time // Дата создания
}
