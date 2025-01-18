package common

import (
	"context"
	memorystorage "github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
	"time"
)

// StorageInterface определяет общие методы для хранилищ.

type StorageInterface interface {
	Connect(ctx context.Context, driverName string, dsn string) error
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, event *event.Event) error
	GetEvent(ctx context.Context, id int64) (*event.Event, error)
	UpdateEvent(ctx context.Context, event *event.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetAllEvents(ctx context.Context) ([]event.Event, error)
	GetAllEventsForDay(ctx context.Context, day time.Time) ([]event.Event, error)
	GetAllEventsForWeek(ctx context.Context, startDayOfWeek time.Time) ([]event.Event, error)
	GetAllEventsForMonth(ctx context.Context, startDayOfMonth time.Time) ([]event.Event, error)
	CreateNotification(ctx context.Context, notification scheduler.Notification) (scheduler.Notification, error)
	DeleteNotificationByEventId(ctx context.Context, eventId int64) error
}

// Фабрика для создания хранилища.

func NewStorage(ctx context.Context, inMemory bool, driverName string, dsn string) (StorageInterface, error) {
	if inMemory {
		storage := memorystorage.New()
		return storage, nil
	}

	storage := sqlstorage.New()
	if err := storage.Connect(ctx, driverName, dsn); err != nil {
		return nil, err
	}
	return storage, nil
}
