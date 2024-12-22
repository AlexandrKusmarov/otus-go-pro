package scheduler

import "time"

type Notification struct {
	ID            int64     // ID Уведомления
	EventID       int64     // ID события, к которому относится уведомление
	Title         string    // Заголовок события
	EventDateTime time.Time // Дата и время события
	UserID        int64     // ID пользователя, которому отправляется уведомление
}
