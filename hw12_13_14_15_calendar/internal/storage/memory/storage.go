package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
)

type Storage struct {
	events        map[int64]*event.Event
	notifications map[int64]*scheduler.Notification
	mu            sync.RWMutex
}

func (s *Storage) Connect(ctx context.Context, driverName string, dsn string) error {
	println(ctx, driverName, dsn)
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	// Для InMemory ничего не требуется
	fmt.Println("Closing in-memory storage")
	return nil
}

func New() *Storage {
	return &Storage{events: make(map[int64]*event.Event), notifications: make(map[int64]*scheduler.Notification)}
}

// GetEvent возвращает событие по его ID.
func (s *Storage) GetEvent(_ context.Context, id int64) (*event.Event, error) {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	event := s.events[id]
	return event, nil
}

// GetAllEvents возвращает все события.
func (s *Storage) GetAllEvents(_ context.Context) ([]event.Event, error) {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	// Создаем копию карты для безопасного возврата
	eventsCopy := make(map[int64]*event.Event)
	for id, event := range s.events {
		eventsCopy[id] = event
	}

	var values []event.Event
	for _, event := range eventsCopy {
		values = append(values, *event)
	}

	return values, nil
}

// GetNotification возвращает уведомление по его ID.
func (s *Storage) GetNotification(id int64) (*scheduler.Notification, bool) {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	notification, exists := s.notifications[id]
	return notification, exists
}

// GetAllNotifications возвращает все уведомления.
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

// AddEvent добавляет или обновляет событие по его ID.
func (s *Storage) CreateEvent(_ context.Context, event *event.Event) error {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	s.events[event.ID] = event
	return nil
}

// RemoveEvent удаляет событие по его ID.
func (s *Storage) DeleteEvent(_ context.Context, id int64) error {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	delete(s.events, id)
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event *event.Event) error {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	s.events[event.ID] = event

	return nil
}

// AddNotification добавляет или обновляет уведомление по его ID.
func (s *Storage) AddNotification(id int64, notification *scheduler.Notification) {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	s.notifications[id] = notification
}

// RemoveNotification удаляет уведомление по его ID.
func (s *Storage) RemoveNotification(id int64) {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	delete(s.notifications, id)
}
