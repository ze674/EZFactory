package models

import (
	"time"
)

type Task struct {
	ID          int       // Уникальный идентификатор
	ProductID   int       // Связь с продуктом
	ProductName string    // Название продукта (для удобства отображения)
	Date        string    // Дата в формате ДД.ММ.ГГГГ
	BatchNumber string    // Номер партии
	Status      string    // Статус: "новое", "в работе", "завершено"
	CreatedAt   time.Time // Дата создания
}
