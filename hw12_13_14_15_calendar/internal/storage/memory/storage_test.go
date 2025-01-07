package memorystorage

import (
	"testing"

	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
)

func TestStorage(t *testing.T) {
	storage := New()

	// Тестирование добавления и получения события
	event1 := &event.Event{ID: 1, Title: "Test Event 1"}
	storage.CreateEvent(nil, event1)

	getEvent, err := storage.GetEvent(nil, event1.ID)
	if err != nil {
		t.Errorf("expected to get event %v, got %v", getEvent, err)
	}

	// Тестирование получения несуществующего события
	_, err = storage.GetEvent(nil, 999)

	if err != nil {
		t.Error("expected event not to exist")
	}

	// Тестирование получения всех событий
	allEvents, err := storage.GetAllEvents(nil)
	if len(allEvents) != 1 {
		t.Errorf("expected 1 event, got %d", len(allEvents))
	}

	// Тестирование удаления события
	storage.DeleteEvent(nil, event1.ID)
	_, err = storage.GetEvent(nil, event1.ID)

	if err != nil {
		t.Error("expected event to be removed")
	}

	// Тестирование добавления и получения уведомления
	notification1 := &scheduler.Notification{ID: 1, Title: "Test Notification 1"}
	storage.AddNotification(notification1.ID, notification1)

	if n, exists := storage.GetNotification(notification1.ID); !exists || n.Title != notification1.Title {
		t.Errorf("expected to get notification %v, got %v", notification1, n)
	}

	// Тестирование получения несуществующего уведомления
	if _, exists := storage.GetNotification(999); exists {
		t.Error("expected notification not to exist")
	}

	// Тестирование получения всех уведомлений
	allNotifications := storage.GetAllNotifications()
	if len(allNotifications) != 1 {
		t.Errorf("expected 1 notification, got %d", len(allNotifications))
	}

	// Тестирование удаления уведомления
	storage.RemoveNotification(notification1.ID)
	if _, exists := storage.GetNotification(notification1.ID); exists {
		t.Error("expected notification to be removed")
	}
}
