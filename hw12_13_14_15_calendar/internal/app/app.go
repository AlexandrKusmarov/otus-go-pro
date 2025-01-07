package app

import (
	"context"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
)

type App struct { // TODO
	storage common.StorageInterface
}

type Logger interface { // TODO
}

type Storage interface { // TODO

}

func New(logger Logger, storage common.StorageInterface) *App {
	fmt.Println(logger, storage) // антилинтер
	return &App{storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event *event.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) GetEvent(ctx context.Context, id int64) (*event.Event, error) {
	return a.storage.GetEvent(ctx, id)
}

func (a *App) UpdateEvent(ctx context.Context, event *event.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetAllEvents(ctx context.Context) ([]event.Event, error) {
	return a.storage.GetAllEvents(ctx)
}
