package models

import (
	"time"
)

type TaskHistory struct {
	ID        int       // Уникальный идентификатор
	TaskID    int       // ID задания
	OldStatus string    // Предыдущий статус
	NewStatus string    // Новый статус
	ChangedAt time.Time // Время изменения
}
