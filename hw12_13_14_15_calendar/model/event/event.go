package event

import (
	"time"
)

type Event struct {
	ID                int64     // Уникальный идентификатор события
	Title             string    // Заголовок события
	EventDateTime     time.Time // Дата и время события
	EventEndDateTime  time.Time // Длительность события
	Description       string    // Описание события (опционально)
	UserID            int64     // ID пользователя, владельца события
	NotifyBeforeEvent time.Time // За сколько времени высылать уведомление (опционально)
}
