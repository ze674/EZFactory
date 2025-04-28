package models

import (
	"time"
)

// Константы для статусов кодов маркировки
const (
	MarkCodeStatusNew     = "new"     // Новый код (не использовался)
	MarkCodeStatusUsed    = "used"    // Использованный код
	MarkCodeStatusInvalid = "invalid" // Недействительный код
)

type MarkCode struct {
	ID           int       // Уникальный идентификатор
	TaskID       int       // ID связанного задания
	Code         string    // Код маркировки
	Status       string    // Статус кода (new, used, invalid)
	UsedAt       time.Time // Когда код был использован (если использован)
	FilePosition int       // Позиция кода в исходном файле
}
