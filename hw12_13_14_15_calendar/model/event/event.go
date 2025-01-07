package event

import (
	"time"
)

type Event struct {
	ID                int64     `db:"id"`                  // Уникальный идентификатор события
	Title             string    `db:"title"`               // Заголовок события
	EventDateTime     time.Time `db:"event_date_time"`     // Дата и время события
	EventEndDateTime  time.Time `db:"event_end_date_time"` // Длительность события
	Description       string    `db:"description"`         // Описание события (опционально)
	UserID            int64     `db:"user_id"`             // ID пользователя, владельца события
	NotifyBeforeEvent time.Time `db:"notify_before_event"` // За сколько времени высылать уведомление (опционально)
}
