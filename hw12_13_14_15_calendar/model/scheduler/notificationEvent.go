package scheduler

import "time"

type Notification struct {
	ID            string    // ID Уведомления
	EventID       string    // ID события, к которому относится уведомление
	Title         string    // Заголовок события
	EventDateTime time.Time // Дата и время события
	UserID        string    // ID пользователя, которому отправляется уведомление
}
