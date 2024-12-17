package common

import (
	"context"
	memorystorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/sql"
)

// StorageInterface определяет общие методы для хранилищ.

type StorageInterface interface {
	Connect(ctx context.Context, driverName string, dsn string) error
	Close(ctx context.Context) error
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
