package event

import (
	"time"
)

type Event struct {
	ID                string    // Уникальный идентификатор события
	Title             string    // Заголовок события
	EventDateTime     time.Time // Дата и время события
	EventEndDateTime  time.Time // Длительность события
	Description       string    // Описание события (опционально)
	UserID            string    // ID пользователя, владельца события
	NotifyBeforeEvent time.Time // За сколько времени высылать уведомление (опционально)
}
