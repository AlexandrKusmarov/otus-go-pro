package memorystorage

import (
	"context"
	"fmt"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/model/event"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/model/scheduler"
	"sync"
)

type Storage struct {
	events        map[int64]*event.Event
	notifications map[int64]*scheduler.Notification
	mu            sync.RWMutex //nolint:unused
}

func (m *Storage) Connect(ctx context.Context, driverName string, dsn string) error {
	return nil
}

func (m *Storage) Close(ctx context.Context) error {
	// Для InMemory ничего не требуется
	fmt.Println("Closing in-memory storage")
	return nil
}

func New() *Storage {
	return &Storage{events: make(map[int64]*event.Event), notifications: make(map[int64]*scheduler.Notification)}
}

// GetEvent возвращает событие по его ID
func (s *Storage) GetEventById(id int64) (*event.Event, bool) {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	return event, exists
}

// GetAllEvents возвращает все события
func (s *Storage) GetAllEvents() map[int64]*event.Event {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	// Создаем копию карты для безопасного возврата
	eventsCopy := make(map[int64]*event.Event)
	for id, event := range s.events {
		eventsCopy[id] = event
	}
	return eventsCopy
}

// GetNotification возвращает уведомление по его ID
func (s *Storage) GetNotification(id int64) (*scheduler.Notification, bool) {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	notification, exists := s.notifications[id]
	return notification, exists
}

// GetAllNotifications возвращает все уведомления
func (s *Storage) GetAllNotifications() map[int64]*scheduler.Notification {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	// Создаем копию карты для безопасного возврата
	notificationsCopy := make(map[int64]*scheduler.Notification)
	for id, notification := range s.notifications {
		notificationsCopy[id] = notification
	}
	return notificationsCopy
}

// AddEvent добавляет или обновляет событие по его ID
func (s *Storage) AddEvent(id int64, event *event.Event) {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	s.events[id] = event
}

// RemoveEvent удаляет событие по его ID
func (s *Storage) RemoveEvent(id int64) {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	delete(s.events, id)
}

// AddNotification добавляет или обновляет уведомление по его ID
func (s *Storage) AddNotification(id int64, notification *scheduler.Notification) {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	s.notifications[id] = notification
}

// RemoveNotification удаляет уведомление по его ID
func (s *Storage) RemoveNotification(id int64) {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	delete(s.notifications, id)
}
