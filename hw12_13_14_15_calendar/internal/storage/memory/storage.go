package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
)

type Storage struct {
	events        map[int64]*event.Event
	notifications map[int64]*scheduler.Notification
	mu            sync.RWMutex
}

func (s *Storage) GetAllEventsForDay(_ context.Context, day time.Time) ([]event.Event, error) {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	var eventsForDay []event.Event
	startOfDay := day.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	for _, e := range s.events {
		if e.EventDateTime.After(startOfDay) && e.EventDateTime.Before(endOfDay) {
			eventsForDay = append(eventsForDay, *e)
		}
	}

	return eventsForDay, nil
}

func (s *Storage) GetAllEventsForWeek(_ context.Context, startDayOfWeek time.Time) ([]event.Event, error) {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	var eventsForWeek []event.Event
	startOfWeek := startDayOfWeek.Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	for _, e := range s.events {
		if e.EventDateTime.After(startOfWeek) && e.EventDateTime.Before(endOfWeek) {
			eventsForWeek = append(eventsForWeek, *e)
		}
	}

	return eventsForWeek, nil
}

func (s *Storage) GetAllEventsForMonth(_ context.Context, startDayOfMonth time.Time) ([]event.Event, error) {
	s.mu.RLock() // Чтение блокировки
	defer s.mu.RUnlock()

	var eventsForMonth []event.Event
	startOfMonth := time.Date(startDayOfMonth.Year(), startDayOfMonth.Month(), 1, 0, 0, 0, 0, startDayOfMonth.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0) // Переход к следующему месяцу

	for _, e := range s.events {
		if e.EventDateTime.After(startOfMonth) && e.EventDateTime.Before(endOfMonth) {
			eventsForMonth = append(eventsForMonth, *e)
		}
	}

	return eventsForMonth, nil
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
func (s *Storage) CreateNotification(_ context.Context, notification scheduler.Notification) (scheduler.Notification, error) {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	s.notifications[notification.ID] = &notification

	return notification, nil
}

// RemoveNotification удаляет уведомление по его ID.
func (s *Storage) RemoveNotification(id int64) {
	s.mu.Lock() // Запись блокировки
	defer s.mu.Unlock()

	delete(s.notifications, id)
}
